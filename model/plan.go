package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	PlanStatusEnabled  = 1
	PlanStatusDisabled = 0

	UserPlanStatusActive  = 1
	UserPlanStatusExpired = 0

	BillingTypeRequest   = "request"
	BillingTypeToken     = "token"
	BillingTypePerRequest = "per_request"

	WindowTypePeriod = "period"
	WindowTypeWeek   = "week"
	WindowTypeMonth  = "month"
)

type Plan struct {
	Id           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string  `gorm:"type:varchar(100);not null" json:"name"`
	Description  string  `gorm:"type:text" json:"description"`
	Price        float64 `gorm:"type:decimal(10,2);default:0" json:"price"`
	DurationDays int     `gorm:"not null;default:30" json:"duration_days"`
	DurationText string  `gorm:"type:varchar(50)" json:"duration_text"`
	Status       int     `gorm:"not null;default:1" json:"status"`
	Recommended  bool    `gorm:"not null;default:false" json:"recommended"`
	Sort         int     `gorm:"not null;default:0" json:"sort"`
	Features     string  `gorm:"type:text" json:"features"`
	ModelLimits  string  `gorm:"type:text;column:model_limits" json:"model_limits"`
	DefaultModel string  `gorm:"type:varchar(100);default:''" json:"default_model"`
	CreatedTime  int64   `gorm:"not null;default:0" json:"created_time"`
	UpdatedTime  int64   `gorm:"not null;default:0" json:"updated_time"`
	CreatedAt    int64   `json:"created_at" gorm:"bigint;default:0"`
	UpdatedAt    int64   `json:"updated_at" gorm:"bigint;default:0"`
}

type ModelLimitRule struct {
	PeriodH        int   `json:"period_h"`
	RequestPeriod  int64 `json:"request_period"`
	RequestWeek    int64 `json:"request_week"`
	RequestMonth   int64 `json:"request_month"`
	TokenPeriod    int64 `json:"token_period"`
	TokenWeek      int64 `json:"token_week"`
	TokenMonth     int64 `json:"token_month"`
}

func (p *Plan) ValidateDefaultModel() error {
	if p.DefaultModel == "" {
		return nil
	}
	limits := p.GetModelLimits()
	if limits == nil {
		return fmt.Errorf("default_model '%s' is set but model_limits is empty", p.DefaultModel)
	}
	if _, ok := limits[p.DefaultModel]; !ok {
		return fmt.Errorf("default_model '%s' is not found in model_limits", p.DefaultModel)
	}
	return nil
}

func (p *Plan) GetModelLimits() map[string]ModelLimitRule {
	if p.ModelLimits == "" {
		return nil
	}
	var limits map[string]ModelLimitRule
	err := json.Unmarshal([]byte(p.ModelLimits), &limits)
	if err != nil {
		return nil
	}
	return limits
}

func (p *Plan) Insert() error {
	return DB.Create(p).Error
}

func (p *Plan) Update() error {
	return DB.Model(p).Select("name", "description", "price", "duration_days",
		"duration_text", "status", "recommended", "sort", "features",
		"model_limits", "default_model", "updated_time").Updates(p).Error
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
	StartTime   int64  `gorm:"not null" json:"start_time"`
	EndTime     int64  `gorm:"not null;index:idx_end_time" json:"end_time"`
	Status      int    `gorm:"not null;default:1;index:idx_user_status" json:"status"`
	BillingType string `gorm:"type:varchar(20);not null;default:'token'" json:"billing_type"`
	Notes       string `gorm:"type:text" json:"notes"`
	CreatedTime int64  `gorm:"not null;default:0" json:"created_time"`
	UpdatedTime int64  `gorm:"not null;default:0" json:"updated_time"`
	CreatedAt   int64  `json:"created_at" gorm:"bigint;default:0"`
	UpdatedAt   int64  `json:"updated_at" gorm:"bigint;default:0"`

	Plan *Plan `gorm:"-" json:"plan,omitempty"`
}

func (up *UserPlan) Insert() error {
	return DB.Create(up).Error
}

func (up *UserPlan) Update() error {
	return DB.Model(up).Select("start_time", "end_time", "status",
		"billing_type", "notes", "updated_time").Updates(up).Error
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
					"used_requests":      pu.Requests,
					"used_prompt_tokens": pu.PromptTokens,
					"used_completion_tokens": pu.CompletionTokens,
					"used_cached_tokens": pu.CachedTokens,
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
			"billing_type": up.BillingType,
			"usage":        usageMap,
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