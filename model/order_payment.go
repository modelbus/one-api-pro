package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/modelbus/one-api-pro/common/helper"
)

// GenerateOrderNo returns a globally unique order number with the
// prefix "TB" (Subscription) or "UP" (Upgrade) followed by a 14-digit
// UTC timestamp and a 6-digit random suffix. Total length: 22 chars.
//
// The function retries up to 8 times if the generated order_no already
// exists (UNIQUE index on orders.order_no) and finally gives up.
func GenerateOrderNo(prefix string) (string, error) {
	now := time.Now().UTC()
	stamp := now.Format("20060102150405")
	rng := rand.New(rand.NewSource(now.UnixNano()))

	for i := 0; i < 8; i++ {
		suffix := fmt.Sprintf("%06d", rng.Intn(1_000_000))
		candidate := prefix + stamp + suffix
		var existing int
		if err := DB.Model(&Order{}).Select("id").Where("order_no = ?", candidate).Scan(&existing).Error; err != nil {
			return "", err
		}
		if existing == 0 {
			return candidate, nil
		}
	}
	return "", errors.New("生成订单号失败，请重试")
}

// SnapshotPlan serializes a Plan into a JSON string used to store an
// immutable copy on the order row.
func SnapshotPlan(plan *Plan) string {
	if plan == nil {
		return ""
	}
	b, err := json.Marshal(plan)
	if err != nil {
		return ""
	}
	return string(b)
}

// GetPlanByOrderPlanInfo parses a previously SnapshotPlan'd string back
// into a Plan. Returns nil on parse failure or empty input.
func GetPlanByOrderPlanInfo(planInfo string) *Plan {
	if planInfo == "" {
		return nil
	}
	var p Plan
	if err := json.Unmarshal([]byte(planInfo), &p); err != nil {
		return nil
	}
	return &p
}

// CalculateUpgradePrice implements the spec's price-diff formula:
//
//	remaining_days = ceil((old_plan_end_time - now) / 86400)
//	old_daily      = old_plan.price / 30
//	new_daily      = new_plan.price / 30
//	upgrade_price  = max(0, (new_daily - old_daily) * remaining_days)
//
// Returns 0 if old is expired (end_time <= now) or old/new prices are
// not strictly decreasing.
func CalculateUpgradePrice(oldPlan *Plan, oldEndTime, now int64, newPlan *Plan) float64 {
	if oldPlan == nil || newPlan == nil {
		return 0
	}
	if oldEndTime <= now {
		return newPlan.Price
	}
	remaining := int64(math.Ceil(float64(oldEndTime-now) / 86400.0))
	if remaining < 0 {
		remaining = 0
	}
	oldDaily := oldPlan.Price / 30.0
	newDaily := newPlan.Price / 30.0
	diff := (newDaily - oldDaily) * float64(remaining)
	if diff < 0 {
		return 0
	}
	return math.Round(diff*100) / 100
}

// CreatePlanOrderInput is the parameter struct accepted by CreatePlanOrder.
type CreatePlanOrderInput struct {
	UserId    int
	PlanId    int
	PayMethod string
	Source    int // OrderSourceUserSelf or OrderSourceAdmin
	Notes     string
}

// CreatePlanOrderOutput is returned by CreatePlanOrder. It carries
// the persisted Order plus the pre-payment fields needed by the caller
// (amount, pay_url for WeChat/Alipay, package_name for display).
type CreatePlanOrderOutput struct {
	Order      *Order
	Mode       string // OrderUpgradeModeStack or OrderUpgradeModePriceDiff (empty if no existing plan)
	Amount     float64
	PackageName string
}

// CreatePlanOrder builds a new order row (type=1) and applies upgrade
// rules. It does NOT call payment pre-pay nor ActivatePackageByOrder;
// those are the caller's responsibility.
func CreatePlanOrder(in CreatePlanOrderInput) (*CreatePlanOrderOutput, error) {
	if in.UserId == 0 {
		return nil, errors.New("user_id 不能为空")
	}
	if in.PlanId == 0 {
		return nil, errors.New("plan_id 不能为空")
	}
	if in.PayMethod == "" {
		return nil, errors.New("pay_method 不能为空")
	}
	if in.Source == 0 {
		in.Source = OrderSourceUserSelf
	}

	plan, err := GetPlanById(in.PlanId)
	if err != nil {
		return nil, fmt.Errorf("套餐不存在: %w", err)
	}
	if plan.Status != PlanStatusEnabled {
		return nil, errors.New("套餐已下架")
	}

	now := helper.GetTimestamp()
	upgradeMode := GetSystemSettingString(SystemSettingKeyPlanUpgradeMode)
	if upgradeMode == "" {
		upgradeMode = OrderUpgradeModePriceDiff
	}

	// Find any current active plan for upgrade logic.
	actives, err := GetActiveUserPlansByUserId(in.UserId)
	if err != nil {
		return nil, err
	}

	amount := plan.Price
	mode := OrderUpgradeModeStack
	if len(actives) > 0 {
		currentPlan := actives[0]
		var currentPlanInfo *Plan
		if currentPlan.Plan != nil {
			currentPlanInfo = currentPlan.Plan
		}
		if upgradeMode == OrderUpgradeModeStack {
			// Stack mode: keep current plan untouched, just create a new one.
			mode = OrderUpgradeModeStack
			amount = plan.Price
		} else {
			// price_diff: only allow upgrade to strictly higher sort.
			if currentPlanInfo != nil && plan.Sort > currentPlanInfo.Sort {
				mode = OrderUpgradeModePriceDiff
				amount = CalculateUpgradePrice(currentPlanInfo, currentPlan.EndTime, now, plan)
			} else if currentPlanInfo != nil && plan.Sort == currentPlanInfo.Sort {
				return nil, errors.New("您已经订阅了同级别的套餐")
			} else {
				return nil, errors.New("不能降级到低级别套餐")
			}
		}
	}

	prefix := "TB"
	if mode == OrderUpgradeModePriceDiff {
		prefix = "UP"
	}
	orderNo, err := GenerateOrderNo(prefix)
	if err != nil {
		return nil, err
	}

	order := &Order{
		Type:       OrderTypePlanSubscription,
		Source:     in.Source,
		OrderNo:    orderNo,
		UserId:     in.UserId,
		PlanId:     in.PlanId,
		PlanInfo:   SnapshotPlan(plan),
		Amount:     amount,
		Status:     OrderStatusPending,
		PayStatus:  OrderPayStatusPending,
		PayMethod:  in.PayMethod,
		CreateTime: now,
		UpdateTime: now,
	}
	if err := order.Insert(); err != nil {
		return nil, err
	}

	return &CreatePlanOrderOutput{
		Order:       order,
		Mode:        mode,
		Amount:      amount,
		PackageName: plan.Name,
	}, nil
}

// ActivatePackageByOrder marks the order as paid and grants the
// corresponding subscription to the user.
//
// mode must be OrderUpgradeModeStack or OrderUpgradeModePriceDiff.
//   - stack: existing active user_plans are kept untouched.
//   - price_diff: existing active user_plans are deactivated before a
//     new one is inserted (upgrade semantics).
//
// On success, the order is updated to status=1, pay_status=1, and
// pay_time is set. A new user_plans row is inserted.
//
// This function is intentionally idempotent at the order level: if the
// order is already paid, it returns nil without re-activating.
func ActivatePackageByOrder(order *Order, mode string) error {
	if order == nil {
		return errors.New("order 不能为空")
	}
	if order.Status == OrderStatusPaid {
		return nil // already activated
	}
	plan := GetPlanByOrderPlanInfo(order.PlanInfo)
	if plan == nil {
		// Fall back to the live plan row if snapshot is missing/corrupt.
		var err error
		plan, err = GetPlanById(order.PlanId)
		if err != nil {
			return fmt.Errorf("套餐不存在: %w", err)
		}
	}

	now := helper.GetTimestamp()
	endTime := now + int64(plan.DurationDays)*86400

	if mode == OrderUpgradeModePriceDiff {
		// Deactivate all currently active user_plans for this user.
		if err := DB.Model(&UserPlan{}).
			Where("user_id = ? AND status = ?", order.UserId, UserPlanStatusActive).
			Updates(map[string]interface{}{
				"status":       UserPlanStatusExpired,
				"updated_time": now,
			}).Error; err != nil {
			return err
		}
	}

	up := &UserPlan{
		UserId:      order.UserId,
		PlanId:      order.PlanId,
		OrderId:     order.Id,
		StartTime:   now,
		EndTime:     endTime,
		Status:      UserPlanStatusActive,
		BillingType: BillingTypeToken,
		CreatedTime: now,
		UpdatedTime: now,
	}
	if err := up.Insert(); err != nil {
		return err
	}

	if err := order.MarkOrderPaid(order.PayMethod, order.PayTradeNo); err != nil {
		return err
	}
	CacheDeleteUserActivePlans(order.UserId)
	return nil
}

// CancelOrder marks a pending order as canceled. Only orders with
// status=0 (pending) can be canceled.
func CancelOrder(order *Order) error {
	if order.Status != OrderStatusPending {
		return errors.New("只能取消待支付订单")
	}
	now := helper.GetTimestamp()
	order.Status = OrderStatusCanceled
	order.PayStatus = OrderPayStatusPending
	order.UpdateTime = now
	return DB.Model(order).Select("status", "pay_status", "update_time").Updates(order).Error
}

// MarkOrderRefunded sets status=3 and pay_status=-1. Does not touch
// the user_plans table (refund flow: depends on actual payment provider
// refund API; not implemented in this milestone).
func MarkOrderRefunded(order *Order) error {
	now := helper.GetTimestamp()
	order.Status = OrderStatusRefunded
	order.PayStatus = OrderPayStatusRefunded
	order.UpdateTime = now
	return DB.Model(order).Select("status", "pay_status", "update_time").Updates(order).Error
}

// IsValidPayMethod returns true if the value is one of the supported
// pay_method strings.
func IsValidPayMethod(m string) bool {
	switch strings.ToLower(m) {
	case OrderPayMethodWechat, OrderPayMethodAlipay,
		OrderPayMethodBank, OrderPayMethodOffline, OrderPayMethodFree:
		return true
	}
	return false
}
