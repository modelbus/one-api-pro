package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StringSlice 在 JSON wire 格式与 DB text 列之间作为 []string 双向桥接。
//
// 写入方向：
//   - MarshalJSON → ["A","B","C"]   （前端 API 看到的形态）
//   - Value       → JSON 字符串     （DB 列仍是 text，无需迁移）
//
// 读取方向（兼容历史三种形态，零迁移）：
//   - 已是数组                → 原样解析
//   - JSON 字符串             ["A","B","C"] → 数组
//   - 换行分隔的纯文本        "A\nB\nC"     → 按行拆分
//   - JSON 对象               {"A":true,...} → 取键名
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
type StringSlice []string

// MarshalJSON 输出 JSON 数组；nil 时输出 [] 而非 null。
//
// 版本: v0.0.12
// 日期: 2026-09-07
func (s StringSlice) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(s))
}

// UnmarshalJSON 接受 JSON 数组；同时兼容旧版字符串/换行文本。
//
// 版本: v0.0.12
// 日期: 2026-09-07
func (s *StringSlice) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*s = nil
		return nil
	}
	// 主路径：JSON 数组
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*s = arr
		return nil
	}
	// 兼容：旧版字符串（可能为 JSON 字符串或换行文本）
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("StringSlice: cannot unmarshal %s: %w", string(data), err)
	}
	*s = splitNonEmptyLines(str)
	return nil
}

// Value 实现 driver.Valuer：DB 写入时序列化为 JSON 字符串。
//
// 版本: v0.0.12
// 日期: 2026-09-07
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	b, err := json.Marshal([]string(s))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan 实现 sql.Scanner：DB 读出时按 JSON 数组解析；失败时按换行兜底。
//
// 版本: v0.0.12
// 日期: 2026-09-07
func (s *StringSlice) Scan(src interface{}) error {
	if src == nil {
		*s = nil
		return nil
	}
	var raw []byte
	switch v := src.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return fmt.Errorf("StringSlice: cannot scan %T", src)
	}
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		*s = nil
		return nil
	}
	// 主路径：JSON 数组
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		*s = arr
		return nil
	}
	// 兜底：旧版换行分隔文本
	*s = splitNonEmptyLines(string(raw))
	return nil
}

// splitNonEmptyLines 按换行拆分并 trim；空行跳过。
//
// 版本: v0.0.12
// 日期: 2026-09-07
func splitNonEmptyLines(s string) []string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

const (
	PlanStatusEnabled  = 1
	PlanStatusDisabled = 0

	UserPlanStatusActive  = 1
	UserPlanStatusExpired = 0

	BillingTypeRequest    = "request"
	BillingTypeToken      = "token"
	BillingTypePerRequest = "per_request"

	WindowTypePeriod = "period"
	WindowTypeWeek   = "week"
	WindowTypeMonth  = "month"
)

type Plan struct {
	Id           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string      `gorm:"type:varchar(100);not null" json:"name"`
	Description  string      `gorm:"type:text" json:"description"`
	Price        float64     `gorm:"type:decimal(10,2);default:0" json:"price"`
	DurationDays int         `gorm:"not null;default:30" json:"duration_days"`
	DurationText string      `gorm:"type:varchar(50)" json:"duration_text"`
	Status       int         `gorm:"not null;default:1" json:"status"`
	Recommended  bool        `gorm:"not null;default:false" json:"recommended"`
	Sort         int         `gorm:"not null;default:0" json:"sort"`
	Features     StringSlice `gorm:"type:text" json:"features"`
	ModelLimits  string      `gorm:"type:text;column:model_limits" json:"model_limits"`
	// BillingType 套餐默认计费维度（token / request）。历史实现缺失该字段，导致
	// user_plans.billing_type 在激活时被硬编码为 token，管理员在前端选择的
	// 「按请求次数」从未生效。存储层单位见 controller/quota_dto.go。
	BillingType string `gorm:"type:varchar(20);not null;default:'token'" json:"billing_type"`
	// VirtualAmount 套餐虚拟余额总池（微元，0 = 不限额度）。
	// 只要套餐未过期，配置了限额的模型都从这笔共享余额里扣；余额耗尽后整个套餐
	// 不可用，请求回落到用户全局余额按量计费。
	VirtualAmount int64 `gorm:"not null;default:0" json:"virtual_amount"`
	CreatedTime   int64 `gorm:"not null;default:0" json:"created_time"`
	UpdatedTime   int64 `gorm:"not null;default:0" json:"updated_time"`
	CreatedAt     int64 `json:"created_at" gorm:"bigint;default:0"`
	UpdatedAt     int64 `json:"updated_at" gorm:"bigint;default:0"`
}

type ModelLimitRule struct {
	PeriodH       int   `json:"period_h"`
	RequestPeriod int64 `json:"request_period"`
	RequestWeek   int64 `json:"request_week"`
	RequestMonth  int64 `json:"request_month"`
	TokenPeriod   int64 `json:"token_period"`
	TokenWeek     int64 `json:"token_week"`
	TokenMonth    int64 `json:"token_month"`
}

// ParseModelLimits 解析 model_limits JSON。
// 空字符串返回 (nil, nil)；JSON 非法时返回错误，供保存前拦截脏配置。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func (p *Plan) ParseModelLimits() (map[string]ModelLimitRule, error) {
	if strings.TrimSpace(p.ModelLimits) == "" {
		return nil, nil
	}
	var limits map[string]ModelLimitRule
	if err := json.Unmarshal([]byte(p.ModelLimits), &limits); err != nil {
		return nil, err
	}
	return limits, nil
}

// GetModelLimits 运行时宽松解析：非法 JSON 返回 nil。
// 注意：调用方必须把 nil 当作「配置不可用」而不是「不限量」，否则会重演
// 「配置写错 = 全免费」的漏洞。
func (p *Plan) GetModelLimits() map[string]ModelLimitRule {
	limits, err := p.ParseModelLimits()
	if err != nil {
		return nil
	}
	return limits
}

// ValidateConfig 落库前校验套餐配置的合法性与完整性。
//
// 规则：
//  1. billing_type 必须是 token 或 request；
//  2. virtual_amount 不能为负；
//  3. model_limits 必须非空且可解析为 map[string]ModelLimitRule
//     （格式写错会让 GetModelLimits 返回 nil，历史实现据此把套餐当成「不限量」
//     免费放行，所以必须在保存阶段直接拒绝）；
//  4. 每条规则至少要配置一个与 billing_type 匹配的限额，否则该模型在加权池里
//     权重恒为 0，等于免费。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func (p *Plan) ValidateConfig() error {
	if !IsValidBillingType(p.BillingType) {
		return fmt.Errorf("计费维度只能是 %s 或 %s", BillingTypeToken, BillingTypeRequest)
	}
	if p.VirtualAmount < 0 {
		return errors.New("套餐额度不能为负数")
	}
	if strings.TrimSpace(p.ModelLimits) == "" {
		return errors.New("模型限制不能为空：请至少配置一个模型及其限额")
	}
	limits, err := p.ParseModelLimits()
	if err != nil {
		return fmt.Errorf("模型限制格式错误: %v", err)
	}
	if len(limits) == 0 {
		return errors.New("模型限制不能为空对象")
	}
	for name, rule := range limits {
		if strings.TrimSpace(name) == "" {
			return errors.New("模型限制中存在空的模型名")
		}
		var matched bool
		if p.BillingType == BillingTypeRequest {
			matched = rule.RequestPeriod > 0 || rule.RequestWeek > 0 || rule.RequestMonth > 0
		} else {
			matched = rule.TokenPeriod > 0 || rule.TokenWeek > 0 || rule.TokenMonth > 0
		}
		if !matched {
			return fmt.Errorf("模型 %s 缺少与计费维度 %s 匹配的限额", name, p.BillingType)
		}
	}
	return nil
}

// GetFeatures 返回当前 features 列表（拷贝，调用方修改不会影响原值）。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func (p *Plan) GetFeatures() []string {
	if len(p.Features) == 0 {
		return []string{}
	}
	out := make([]string, len(p.Features))
	copy(out, p.Features)
	return out
}

func (p *Plan) Insert() error {
	return DB.Create(p).Error
}

func (p *Plan) Update() error {
	return DB.Model(p).Select("name", "description", "price", "duration_days",
		"duration_text", "status", "recommended", "sort", "features",
		"model_limits", "billing_type", "virtual_amount", "updated_time").Updates(p).Error
}

func DeletePlanById(id int) error {
	// 先查询再删除，确保 GORM 回调能取到正确的 Id
	var plan Plan
	if err := DB.First(&plan, "id = ?", id).Error; err != nil {
		return err
	}
	return DB.Delete(&plan).Error
}

func GetPlanById(id int) (*Plan, error) {
	var plan Plan
	err := DB.First(&plan, "id = ?", id).Error
	return &plan, err
}

func GetAllPlans(startIdx int, num int) ([]*Plan, error) {
	var plans []*Plan
	err := DB.Order("sort asc, id desc").Limit(num).Offset(startIdx).Find(&plans).Error
	return plans, err
}

func SearchPlans(keyword string) ([]*Plan, error) {
	var plans []*Plan
	keywordCol := "`name`"
	if common.UsingPostgreSQL {
		keywordCol = `"name"`
	}
	err := DB.Where(keywordCol+" LIKE ?", keyword+"%").Order("sort asc, id desc").Find(&plans).Error
	return plans, err
}

type UserPlan struct {
	Id          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      int    `gorm:"not null;index:idx_user_status" json:"user_id"`
	PlanId      int    `gorm:"not null" json:"plan_id"`
	OrderId     int    `gorm:"not null;default:0" json:"order_id"`
	StartTime   int64  `gorm:"not null" json:"start_time"`
	EndTime     int64  `gorm:"not null;index:idx_end_time" json:"end_time"`
	Status      int    `gorm:"not null;default:1;index:idx_user_status" json:"status"`
	BillingType string `gorm:"type:varchar(20);not null;default:'token'" json:"billing_type"`
	// 以下三个字段是「购买时的套餐快照」：plans 表行被删除后，订阅仍能按快照继续
	// 计费，不会因为拿不到套餐配置而退化成「不限量免费」。
	// ModelLimits 快照仅在 plans 行缺失时生效，plans 行存在时以实时配置为准。
	ModelLimits   string `gorm:"type:text;column:model_limits" json:"model_limits"`
	VirtualAmount int64  `gorm:"not null;default:0" json:"virtual_amount"`
	// UsedAmount 订阅期内累计已用额度（微元），不随窗口重置。
	UsedAmount  int64  `gorm:"not null;default:0" json:"used_amount"`
	Notes       string `gorm:"type:text" json:"notes"`
	CreatedTime int64  `gorm:"not null;default:0" json:"created_time"`
	UpdatedTime int64  `gorm:"not null;default:0" json:"updated_time"`
	CreatedAt   int64  `json:"created_at" gorm:"bigint;default:0"`
	UpdatedAt   int64  `json:"updated_at" gorm:"bigint;default:0"`

	Plan *Plan `gorm:"-" json:"plan,omitempty"`
}

// EffectiveModelLimits 返回订阅当前生效的模型限制。
// plans 行仍在时用实时配置；套餐已被删除时回退到购买快照。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func (up *UserPlan) EffectiveModelLimits() (map[string]ModelLimitRule, error) {
	if up.Plan != nil {
		return up.Plan.ParseModelLimits()
	}
	if strings.TrimSpace(up.ModelLimits) == "" {
		return nil, nil
	}
	var limits map[string]ModelLimitRule
	if err := json.Unmarshal([]byte(up.ModelLimits), &limits); err != nil {
		return nil, err
	}
	return limits, nil
}

// EffectiveBillingType 返回订阅当前生效的计费维度。
// plans 行仍在时以套餐配置为准，否则用订阅自身记录（激活时即已写入）。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func (up *UserPlan) EffectiveBillingType() string {
	if up.Plan != nil && IsValidBillingType(up.Plan.BillingType) {
		return up.Plan.BillingType
	}
	if IsValidBillingType(up.BillingType) {
		return up.BillingType
	}
	return BillingTypeToken
}

// EffectiveVirtualAmount 返回订阅当前生效的虚拟余额总池（微元）。
// plans 行仍在时以套餐配置为准，否则回退到购买快照。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func (up *UserPlan) EffectiveVirtualAmount() int64 {
	if up.Plan != nil {
		return up.Plan.VirtualAmount
	}
	return up.VirtualAmount
}

// IsValidBillingType 判断计费维度是否合法（套餐仅支持 token / request）。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func IsValidBillingType(t string) bool {
	return t == BillingTypeToken || t == BillingTypeRequest
}

// GetUserPlanUsedAmount 实时读取订阅已用额度（微元）。
// 不走 Redis 缓存：虚拟余额的耗尽判定必须基于最新值，否则会出现超额放行。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func GetUserPlanUsedAmount(userPlanId int) (int64, error) {
	var up UserPlan
	if err := DB.Select("used_amount").First(&up, "id = ?", userPlanId).Error; err != nil {
		return 0, err
	}
	return up.UsedAmount, nil
}

// IncrementUserPlanUsedAmount 原子累加订阅已用额度（微元）。
// 用 gorm.Expr 直接自增，避免读改写竞态。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func IncrementUserPlanUsedAmount(userPlanId int, quota int64) error {
	if userPlanId <= 0 || quota == 0 {
		return nil
	}
	return DB.Model(&UserPlan{}).Where("id = ?", userPlanId).
		UpdateColumn("used_amount", gorm.Expr("used_amount + ?", quota)).Error
}

// RemainingVirtualAmount 计算虚拟余额剩余额度（微元）。
// virtualAmount <= 0 表示「不限额度」，返回 -1 以便调用方区分「不限」与「已耗尽(0)」。
//
// 版本: v0.0.25
// 日期: 2026-10-05
// 作者: opencode
func RemainingVirtualAmount(virtualAmount, usedAmount int64) int64 {
	if virtualAmount <= 0 {
		return -1
	}
	if usedAmount >= virtualAmount {
		return 0
	}
	return virtualAmount - usedAmount
}

func (up *UserPlan) Insert() error {
	// 防御：plan_id=0 的 user_plan 行无法关联到套餐，会让 admin 列表
	// 把套餐列显示成空白。此处前置校验，避免写入脏数据。
	// Defense: block inserts where plan_id=0 (these rows render as empty
	//套餐 in admin lists and cannot be traced back to a套餐).
	if up.PlanId <= 0 {
		return errors.New("user_plan.plan_id 不能为空")
	}
	return DB.Create(up).Error
}

func (up *UserPlan) Update() error {
	return DB.Model(up).Select("start_time", "end_time", "status",
		"order_id", "billing_type", "notes", "updated_time").Updates(up).Error
}

func DeleteUserPlanById(id int) error {
	// 先查询再删除，确保 GORM 回调（如 cluster 的 beforeDelete）能取到正确的 Id
	var up UserPlan
	if err := DB.First(&up, "id = ?", id).Error; err != nil {
		return err
	}
	return DB.Delete(&up).Error
}

func GetUserPlanById(id int) (*UserPlan, error) {
	var up UserPlan
	err := DB.First(&up, "id = ?", id).Error
	return &up, err
}

// GetUserPlanByOrderId returns the user_plan row that was created for
// the given order id. Returns nil if not found.
func GetUserPlanByOrderId(orderId int) (*UserPlan, error) {
	if orderId == 0 {
		return nil, nil
	}
	var up UserPlan
	if err := DB.Where("order_id = ?", orderId).First(&up).Error; err != nil {
		return nil, err
	}
	return &up, nil
}

func GetActiveUserPlansByUserId(userId int) ([]*UserPlan, error) {
	var ups []*UserPlan
	now := helper.GetTimestamp()
	err := DB.Where("user_id = ? AND status = ? AND end_time > ?", userId, UserPlanStatusActive, now).
		Order("end_time asc").Find(&ups).Error
	if err != nil {
		return nil, err
	}
	for _, up := range ups {
		plan, err := GetPlanById(up.PlanId)
		if err == nil {
			up.Plan = plan
		}
	}
	return ups, nil
}

func GetAllUserPlans(startIdx int, num int, userId int, status int) ([]*UserPlan, error) {
	var ups []*UserPlan
	query := DB.Model(&UserPlan{})
	if userId > 0 {
		query = query.Where("user_id = ?", userId)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Order("id desc").Limit(num).Offset(startIdx).Find(&ups).Error
	if err != nil {
		return nil, err
	}
	for _, up := range ups {
		plan, err := GetPlanById(up.PlanId)
		if err == nil {
			up.Plan = plan
		}
	}
	return ups, nil
}

func SearchUserPlans(keyword string, startIdx int, num int) ([]*UserPlan, error) {
	var ups []*UserPlan
	err := DB.Joins("JOIN users ON users.id = user_plans.user_id").
		Where("users.username LIKE ?", keyword+"%").
		Order("user_plans.id desc").Limit(num).Offset(startIdx).
		Find(&ups).Error
	return ups, err
}

func ExpireUserPlans() {
	now := helper.GetTimestamp()
	result := DB.Model(&UserPlan{}).Where("status = ? AND end_time <= ?", UserPlanStatusActive, now).
		Update("status", UserPlanStatusExpired)
	if result.RowsAffected > 0 {
		DB.RowsAffected = result.RowsAffected
	}
}

type PlanUsage struct {
	Id               uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserPlanId       int    `gorm:"not null;uniqueIndex:idx_plan_unique_window" json:"user_plan_id"`
	Model            string `gorm:"type:varchar(100);not null;uniqueIndex:idx_plan_unique_window" json:"model"`
	WindowType       string `gorm:"type:varchar(20);not null;uniqueIndex:idx_plan_unique_window" json:"window_type"`
	WindowIndex      int    `gorm:"not null;uniqueIndex:idx_plan_unique_window" json:"window_index"`
	Requests         int64  `gorm:"not null;default:0" json:"requests"`
	PromptTokens     int64  `gorm:"not null;default:0" json:"prompt_tokens"`
	CompletionTokens int64  `gorm:"not null;default:0" json:"completion_tokens"`
	CachedTokens     int64  `gorm:"not null;default:0" json:"cached_tokens"`
	UpdatedTime      int64  `gorm:"not null;default:0" json:"updated_time"`
	CreatedAt        int64  `json:"created_at" gorm:"bigint;default:0"`
	UpdatedAt        int64  `json:"updated_at" gorm:"bigint;default:0"`
}

func CalcWindowIndex(now int64, startTime int64, windowType string, periodH int) int {
	elapsed := now - startTime
	switch windowType {
	case WindowTypePeriod:
		if periodH <= 0 {
			periodH = 5
		}
		return int(elapsed / int64(periodH*3600))
	case WindowTypeWeek:
		return int(elapsed / (7 * 86400))
	case WindowTypeMonth:
		return int(elapsed / (30 * 86400))
	}
	return 0
}

func GetWindowDurationSeconds(windowType string, periodH int) int64 {
	switch windowType {
	case WindowTypePeriod:
		if periodH <= 0 {
			periodH = 5
		}
		return int64(periodH) * 3600
	case WindowTypeWeek:
		return 7 * 86400
	case WindowTypeMonth:
		return 30 * 86400
	}
	return 0
}

func GetPlanUsage(userPlanId int, model string, windowType string, windowIndex int) (*PlanUsage, error) {
	var pu PlanUsage
	err := DB.Where("user_plan_id = ? AND model = ? AND window_type = ? AND window_index = ?",
		userPlanId, model, windowType, windowIndex).First(&pu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &PlanUsage{
				UserPlanId:       userPlanId,
				Model:            model,
				WindowType:       windowType,
				WindowIndex:      windowIndex,
				Requests:         0,
				PromptTokens:     0,
				CompletionTokens: 0,
				CachedTokens:     0,
			}, nil
		}
		return nil, err
	}
	return &pu, nil
}

func IncrementPlanUsage(userPlanId int, model string, windowType string, windowIndex int, requests int64, promptTokens int64, completionTokens int64, cachedTokens int64) error {
	now := helper.GetTimestamp()
	if common.UsingPostgreSQL {
		result := DB.Model(&PlanUsage{}).
			Where("user_plan_id = ? AND model = ? AND window_type = ? AND window_index = ?",
				userPlanId, model, windowType, windowIndex).
			Updates(map[string]interface{}{
				"requests":          gorm.Expr("requests + ?", requests),
				"prompt_tokens":     gorm.Expr("prompt_tokens + ?", promptTokens),
				"completion_tokens": gorm.Expr("completion_tokens + ?", completionTokens),
				"cached_tokens":     gorm.Expr("cached_tokens + ?", cachedTokens),
				"updated_time":      now,
			})
		if result.RowsAffected == 0 {
			pu := PlanUsage{
				UserPlanId:       userPlanId,
				Model:            model,
				WindowType:       windowType,
				WindowIndex:      windowIndex,
				Requests:         requests,
				PromptTokens:     promptTokens,
				CompletionTokens: completionTokens,
				CachedTokens:     cachedTokens,
				UpdatedTime:      now,
			}
			return DB.Create(&pu).Error
		}
		return result.Error
	}
	if common.UsingSQLite {
		pu := PlanUsage{
			UserPlanId:       userPlanId,
			Model:            model,
			WindowType:       windowType,
			WindowIndex:      windowIndex,
			Requests:         requests,
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			CachedTokens:     cachedTokens,
			UpdatedTime:      now,
		}
		return DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_plan_id"}, {Name: "model"}, {Name: "window_type"}, {Name: "window_index"}},
			DoUpdates: clause.AssignmentColumns([]string{"requests", "prompt_tokens", "completion_tokens", "cached_tokens", "updated_time"}),
		}).Create(&pu).Error
	}
	sql := fmt.Sprintf(
		"INSERT INTO plan_usages (user_plan_id, model, window_type, window_index, requests, prompt_tokens, completion_tokens, cached_tokens, updated_time) "+
			"VALUES (%d, '%s', '%s', %d, %d, %d, %d, %d, %d) "+
			"ON DUPLICATE KEY UPDATE requests = requests + %d, prompt_tokens = prompt_tokens + %d, completion_tokens = completion_tokens + %d, cached_tokens = cached_tokens + %d, updated_time = %d",
		userPlanId, model, windowType, windowIndex, requests, promptTokens, completionTokens, cachedTokens, now,
		requests, promptTokens, completionTokens, cachedTokens, now,
	)
	return DB.Exec(sql).Error
}

func GetPlanUsageByUserPlanId(userPlanId int) ([]*PlanUsage, error) {
	var pus []*PlanUsage
	err := DB.Where("user_plan_id = ?", userPlanId).Find(&pus).Error
	return pus, err
}

func CleanOldPlanUsage(beforeWindowIndex int, windowType string) error {
	return DB.Where("window_index < ? AND window_type = ?", beforeWindowIndex, windowType).Delete(&PlanUsage{}).Error
}

func GetUserSubscriptionInfo(userId int) ([]map[string]interface{}, error) {
	ups, err := GetActiveUserPlansByUserId(userId)
	if err != nil {
		return nil, err
	}
	now := helper.GetTimestamp()
	var result []map[string]interface{}
	for _, up := range ups {
		if up.Plan == nil {
			continue
		}
		limits := up.Plan.GetModelLimits()
		usageMap := make(map[string]map[string]interface{})
		for model, rule := range limits {
			modelUsage := make(map[string]interface{})
			for _, windowType := range []string{WindowTypePeriod, WindowTypeWeek, WindowTypeMonth} {
				windowIndex := CalcWindowIndex(now, up.StartTime, windowType, rule.PeriodH)
				pu, _ := GetPlanUsage(int(up.Id), model, windowType, windowIndex)
				entry := map[string]interface{}{
					"used_requests":          pu.Requests,
					"used_prompt_tokens":     pu.PromptTokens,
					"used_completion_tokens": pu.CompletionTokens,
					"used_cached_tokens":     pu.CachedTokens,
				}
				switch windowType {
				case WindowTypePeriod:
					entry["limit_requests"] = rule.RequestPeriod
					entry["limit_tokens"] = rule.TokenPeriod
				case WindowTypeWeek:
					entry["limit_requests"] = rule.RequestWeek
					entry["limit_tokens"] = rule.TokenWeek
				case WindowTypeMonth:
					entry["limit_requests"] = rule.RequestMonth
					entry["limit_tokens"] = rule.TokenMonth
				}
				modelUsage[windowType] = entry
			}
			usageMap[model] = modelUsage
		}
		entry := map[string]interface{}{
			"id":           up.Id,
			"plan_id":      up.PlanId,
			"plan_name":    up.Plan.Name,
			"start_time":   up.StartTime,
			"end_time":     up.EndTime,
			"status":       up.Status,
			"billing_type": up.EffectiveBillingType(),
			"usage":        usageMap,
			// 虚拟余额（微元）：由 controller 层换算为「元」。-1 表示不限额度。
			"virtual_amount":   up.EffectiveVirtualAmount(),
			"used_amount":      up.UsedAmount,
			"remaining_amount": RemainingVirtualAmount(up.EffectiveVirtualAmount(), up.UsedAmount),
		}
		result = append(result, entry)
	}
	return result, nil
}

// CacheGetUserActivePlans returns active user plans from cache or DB
var UserPlanCacheSeconds = 300

func CacheGetUserActivePlans(userId int) ([]*UserPlan, error) {
	if !common.RedisEnabled {
		return GetActiveUserPlansByUserId(userId)
	}
	key := fmt.Sprintf("user_plans:%d", userId)
	data, err := common.RedisGet(key)
	if err == nil && data != "" {
		var ups []*UserPlan
		if jsonErr := json.Unmarshal([]byte(data), &ups); jsonErr == nil {
			planMap := make(map[int]*Plan)
			for _, up := range ups {
				if up.Plan == nil {
					if p, ok := planMap[up.PlanId]; ok {
						up.Plan = p
					} else {
						p, err := GetPlanById(up.PlanId)
						if err == nil {
							up.Plan = p
							planMap[up.PlanId] = p
						}
					}
				}
			}
			return ups, nil
		}
	}
	ups, err := GetActiveUserPlansByUserId(userId)
	if err != nil {
		return nil, err
	}
	jsonBytes, _ := json.Marshal(ups)
	_ = common.RedisSet(key, string(jsonBytes), time.Duration(UserPlanCacheSeconds)*time.Second)
	return ups, nil
}

func CacheDeleteUserActivePlans(userId int) {
	if common.RedisEnabled {
		_ = common.RedisDel(fmt.Sprintf("user_plans:%d", userId))
	}
}

type UserSubscriptionBrief struct {
	PlanName    string `json:"plan_name"`
	BillingType string `json:"billing_type"`
	EndTime     int64  `json:"end_time"`
	Status      int    `json:"status"`
}

func GetUserSubscriptionBriefs(userIds []int) (map[int][]*UserSubscriptionBrief, error) {
	var ups []*UserPlan
	err := DB.Where("user_id IN ? AND status = ?", userIds, UserPlanStatusActive).Find(&ups).Error
	if err != nil {
		return nil, err
	}
	planIds := make(map[int]bool)
	for _, up := range ups {
		planIds[up.PlanId] = true
	}
	plans := make(map[int]*Plan)
	for pid := range planIds {
		p, err := GetPlanById(pid)
		if err == nil {
			plans[pid] = p
		}
	}
	result := make(map[int][]*UserSubscriptionBrief)
	for _, up := range ups {
		brief := &UserSubscriptionBrief{
			EndTime:     up.EndTime,
			Status:      up.Status,
			BillingType: up.BillingType,
		}
		if p, ok := plans[up.PlanId]; ok {
			brief.PlanName = p.Name
		}
		result[up.UserId] = append(result[up.UserId], brief)
	}
	return result, nil
}
