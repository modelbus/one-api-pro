package model

import (
	"errors"

	"github.com/modelbus/one-api-pro/common/helper"
)

// Order type constants
const (
	OrderTypePlanSubscription = 1 // 套餐订阅订单
	OrderTypeTopup            = 2 // 充值订单
)

// Order source constants
const (
	OrderSourceUserSelf = 1 // 用户自助下单
	OrderSourceAdmin    = 2 // 管理员下单
)

// Order status constants
const (
	OrderStatusPending  = 0 // 待支付
	OrderStatusPaid     = 1 // 已支付
	OrderStatusCanceled = 2 // 已取消
	OrderStatusRefunded = 3 // 已退款
)

// Order pay_status constants (mirrors status for finer-grained reporting)
const (
	OrderPayStatusRefunded = -1 // 已退款
	OrderPayStatusPending  = 0  // 未支付
	OrderPayStatusPaid     = 1  // 已支付
)

// Order pay_method constants
const (
	OrderPayMethodWechat  = "wechat"
	OrderPayMethodAlipay  = "alipay"
	OrderPayMethodBank    = "bank"      // 预留
	OrderPayMethodOffline = "offline"   // 管理员线下收款
	OrderPayMethodFree    = "free"      // 管理员免费赠送
)

// Order upgrade mode constants
const (
	OrderUpgradeModePriceDiff = "price_diff" // 差价升级（默认）
	OrderUpgradeModeStack      = "stack"      // 叠加
)

type Order struct {
	Id          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Type        int     `gorm:"not null;default:1" json:"type"`
	Source      int     `gorm:"not null;default:1" json:"source"`
	OrderNo     string  `gorm:"type:varchar(64);not null;uniqueIndex:idx_order_no" json:"order_no"`
	UserId      int     `gorm:"not null;index:idx_user_id" json:"user_id"`
	PlanId      int     `gorm:"not null;default:0" json:"plan_id"`
	PlanInfo    string  `gorm:"type:text" json:"plan_info"`
	Amount      float64 `gorm:"type:decimal(10,2)" json:"amount"`
	Status      int     `gorm:"not null;default:0;index:idx_status" json:"status"`
	PayStatus   int     `gorm:"not null;default:0" json:"pay_status"`
	PayMethod   string  `gorm:"type:varchar(20);default:''" json:"pay_method"`
	PayTime     int64   `gorm:"bigint;default:0" json:"pay_time"`
	PayTradeNo  string  `gorm:"type:varchar(64);default:''" json:"pay_trade_no"`
	CreateTime  int64   `gorm:"bigint;default:0;index:idx_create_time" json:"create_time"`
	UpdateTime  int64   `gorm:"bigint;default:0" json:"update_time"`
}

func (o *Order) Insert() error {
	return DB.Create(o).Error
}

func (o *Order) Update() error {
	return DB.Model(o).Select("type", "source", "plan_id", "plan_info", "amount",
		"status", "pay_status", "pay_method", "pay_time", "pay_trade_no",
		"update_time").Updates(o).Error
}

func (o *Order) Delete() error {
	return DB.Delete(o).Error
}

func GetOrderById(id int) (*Order, error) {
	if id == 0 {
		return nil, errors.New("id 为空")
	}
	o := Order{Id: id}
	if err := DB.First(&o, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func GetOrderByOrderNo(orderNo string) (*Order, error) {
	if orderNo == "" {
		return nil, errors.New("order_no 为空")
	}
	o := Order{}
	if err := DB.Where("order_no = ?", orderNo).First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func GetAllOrders(startIdx int, num int, orderType int) ([]*Order, error) {
	var orders []*Order
	q := DB.Order("id desc")
	if orderType > 0 {
		q = q.Where("type = ?", orderType)
	}
	if err := q.Limit(num).Offset(startIdx).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func SearchOrders(keyword string, orderType int) ([]*Order, error) {
	var orders []*Order
	q := DB.Where("order_no LIKE ? OR pay_trade_no LIKE ?", keyword+"%", keyword+"%")
	if orderType > 0 {
		q = q.Where("type = ?", orderType)
	}
	if err := q.Order("id desc").Limit(50).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// GetUserOrders returns orders belonging to userId, newest first.
// orderType=0 means any type. limit=0 means no limit.
func GetUserOrders(userId int, orderType int, limit int) ([]*Order, error) {
	var orders []*Order
	q := DB.Where("user_id = ?", userId)
	if orderType > 0 {
		q = q.Where("type = ?", orderType)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Order("id desc").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// MarkOrderPaid sets status=1, pay_status=1, pay_time, pay_trade_no, update_time.
func (o *Order) MarkOrderPaid(payMethod, payTradeNo string) error {
	o.Status = OrderStatusPaid
	o.PayStatus = OrderPayStatusPaid
	o.PayMethod = payMethod
	o.PayTime = helper.GetTimestamp()
	o.PayTradeNo = payTradeNo
	return DB.Model(o).Select("status", "pay_status", "pay_method", "pay_time", "pay_trade_no", "update_time").Updates(o).Error
}
