// quota_dto.go Controller 层 API 边界的额度单位换算（微元 <-> 元）
// API-boundary quota unit conversion (micro-quota <-> CNY yuan) in controllers
// 版本: v0.0.25
// 日期: 2026-10-04
// 作者: opencode
package controller

import (
	"encoding/json"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/model"
)

// 单位口径约定（v0.0.25 起）：
//   - 存储层（DB / model / relay 计费）：一律「微元」整数， 1 元 = config.QuotaPerUnit = 1_000_000 微元。
//   - API 边界（controller 层）：一律「元」浮点数。出入参在此文件统一换算，
//     model 层与 relay 内部不得感知「元」，前端也不得再做任何 1e6 换算。
//
// 为什么必须放在 controller 层而不是 model 层：
//   - cluster/watcher.go 的 INSERT 事件用 json.Marshal(model) 序列化整行后跨节点推送，
//     若改 model 的 JSON 表示，会把「元」写进对端 bigint 列，造成跨节点数据损坏。
//   - model/cache.go 把 model.Token 的 JSON 存入 Redis，单位变更会让存量缓存被错误解读。
//   - model/order_payment.go 的订单快照同理。
//
// Unit convention (since v0.0.25):
//   - Storage layer (DB / model / relay billing): always integer micro-quota,
//     1 CNY = config.QuotaPerUnit = 1_000_000 micro-quota.
//   - API boundary (controller layer): always CNY yuan as float64. Conversion happens
//     exclusively in this file; model/relay must never see "yuan", and the frontend
//     must never do any 1e6 conversion again.
//
// Why the controller layer instead of the model layer:
//   - cluster/watcher.go serializes a whole row via json.Marshal(model) for INSERT events
//     and pushes it across nodes. Changing the model JSON representation would write
//     "yuan" into the peer's bigint column and corrupt cross-node data.
//   - model/cache.go stores model.Token's JSON in Redis; a unit change would make the
//     existing cache entries be misread.
//   - model/order_payment.go's order snapshot has the same problem.

// quotaToYuan 将存储层的微元额度换算为「元」，用于 API 输出。
// Convert storage-layer micro-quota to CNY yuan for API responses.
//
// 精度说明：quota 上限受 int64 约束，实际业务量级（< 2^53 微元 ≈ 90 亿元）
// 可被 float64 精确表示，故 float64 足够承载展示与再输入。
func quotaToYuan(q int64) float64 {
	return float64(q) / config.QuotaPerUnit
}

// yuanToQuota 将 API 输入的「元」换算为微元额度，用于落库。
// Convert CNY yuan from API input to micro-quota for persistence.
//
// 复用 model.YuanToQuota 以保证与业务层口径完全一致（内部走 math.Round）。
func yuanToQuota(y float64) int64 {
	return model.YuanToQuota(y)
}

// ---------------------------------------------------------------------------
// 输出 DTO（微元 -> 元）
//
// 实现方式：内嵌 model 结构体 + 用同名字段遮蔽（shadow）需要换算的字段。
// encoding/json 选取「层级最浅」的字段，外层 float64 会覆盖内嵌 int64，
// 因此 DTO 只需声明要换算的字段，其余字段自动透传，改动面最小。
//
// Output DTOs (micro-quota -> yuan).
// Implementation: embed the model struct and shadow the amount fields with the same
// JSON name. encoding/json picks the shallowest field, so the outer float64 wins and
// every other field passes through untouched.
// ---------------------------------------------------------------------------

// userDTO 用户信息输出：quota / used_quota 以「元」为单位。
// User response: quota / used_quota expressed in CNY yuan.
type userDTO struct {
	*model.User
	Quota     float64 `json:"quota"`
	UsedQuota float64 `json:"used_quota"`
}

// toUserDTO 把单个用户转换为「元」口径的输出结构。
// Convert a single user to the yuan-denominated response shape.
func toUserDTO(u *model.User) *userDTO {
	if u == nil {
		return nil
	}
	return &userDTO{
		User:      u,
		Quota:     quotaToYuan(u.Quota),
		UsedQuota: quotaToYuan(u.UsedQuota),
	}
}

// toUserDTOs 批量转换用户列表；入参为 nil 时返回空切片而非 nil。
// Convert a user list; returns an empty slice (not nil) when the input is empty.
func toUserDTOs(users []*model.User) []*userDTO {
	out := make([]*userDTO, 0, len(users))
	for _, u := range users {
		out = append(out, toUserDTO(u))
	}
	return out
}

// tokenDTO 令牌输出：remain_quota / used_quota 以「元」为单位。
// Token response: remain_quota / used_quota expressed in CNY yuan.
type tokenDTO struct {
	*model.Token
	RemainQuota float64 `json:"remain_quota"`
	UsedQuota   float64 `json:"used_quota"`
}

// toTokenDTO 把单个令牌转换为「元」口径的输出结构。
// Convert a single token to the yuan-denominated response shape.
func toTokenDTO(t *model.Token) *tokenDTO {
	if t == nil {
		return nil
	}
	return &tokenDTO{
		Token:       t,
		RemainQuota: quotaToYuan(t.RemainQuota),
		UsedQuota:   quotaToYuan(t.UsedQuota),
	}
}

// toTokenDTOs 批量转换令牌列表；入参为 nil 时返回空切片而非 nil。
// Convert a token list; returns an empty slice (not nil) when the input is empty.
func toTokenDTOs(tokens []*model.Token) []*tokenDTO {
	out := make([]*tokenDTO, 0, len(tokens))
	for _, t := range tokens {
		out = append(out, toTokenDTO(t))
	}
	return out
}

// channelDTO 渠道输出：used_quota 以「元」为单位。
// Channel response: used_quota expressed in CNY yuan.
//
// 注意：Channel.Balance（上游账户余额）保持 USD 原样，不做换算，
// 因为它来自各上游厂商的账户币种，系统内并无 CNY/USD 汇率配置。
type channelDTO struct {
	*model.Channel
	UsedQuota float64 `json:"used_quota"`
}

// toChannelDTO 把单个渠道转换为「元」口径的输出结构。
// Convert a single channel to the yuan-denominated response shape.
func toChannelDTO(ch *model.Channel) *channelDTO {
	if ch == nil {
		return nil
	}
	return &channelDTO{
		Channel:   ch,
		UsedQuota: quotaToYuan(ch.UsedQuota),
	}
}

// toChannelDTOs 批量转换渠道列表；入参为 nil 时返回空切片而非 nil。
// Convert a channel list; returns an empty slice (not nil) when the input is empty.
func toChannelDTOs(channels []*model.Channel) []*channelDTO {
	out := make([]*channelDTO, 0, len(channels))
	for _, ch := range channels {
		out = append(out, toChannelDTO(ch))
	}
	return out
}

// redemptionDTO 兑换码输出：quota 以「元」为单位。
// Redemption response: quota expressed in CNY yuan.
type redemptionDTO struct {
	*model.Redemption
	Quota float64 `json:"quota"`
}

// toRedemptionDTO 把单个兑换码转换为「元」口径的输出结构。
// Convert a single redemption code to the yuan-denominated response shape.
func toRedemptionDTO(r *model.Redemption) *redemptionDTO {
	if r == nil {
		return nil
	}
	return &redemptionDTO{
		Redemption: r,
		Quota:      quotaToYuan(r.Quota),
	}
}

// toRedemptionDTOs 批量转换兑换码列表；入参为 nil 时返回空切片而非 nil。
// Convert a redemption list; returns an empty slice (not nil) when the input is empty.
func toRedemptionDTOs(redemptions []*model.Redemption) []*redemptionDTO {
	out := make([]*redemptionDTO, 0, len(redemptions))
	for _, r := range redemptions {
		out = append(out, toRedemptionDTO(r))
	}
	return out
}

// logDTO 日志输出：quota 以「元」为单位。
// Log response: quota expressed in CNY yuan.
type logDTO struct {
	*model.Log
	Quota float64 `json:"quota"`
}

// toLogDTO 把单条日志转换为「元」口径的输出结构。
// Convert a single log entry to the yuan-denominated response shape.
func toLogDTO(l *model.Log) *logDTO {
	if l == nil {
		return nil
	}
	return &logDTO{
		Log:   l,
		Quota: quotaToYuan(int64(l.Quota)),
	}
}

// toLogDTOs 批量转换日志列表；入参为 nil 时返回空切片而非 nil。
// Convert a log list; returns an empty slice (not nil) when the input is empty.
func toLogDTOs(logs []*model.Log) []*logDTO {
	out := make([]*logDTO, 0, len(logs))
	for _, l := range logs {
		out = append(out, toLogDTO(l))
	}
	return out
}

// topupPresetDTO 充值快捷金额输出：bonus_quota 以「元」为单位。
// Top-up preset response: bonus_quota expressed in CNY yuan.
//
// 该结构不是内嵌 DTO，因为 model.TopupPreset 同时被写库路径（system_settings 的 JSON）
// 与 API 复用，而其 DB 与 API 的口径恰好不同（DB 微元 / API 元），无法用同名遮蔽区分。
type topupPresetDTO struct {
	Amount     float64 `json:"amount"`
	BonusQuota float64 `json:"bonus_quota"`
}

// toTopupPresetDTOs 批量转换充值快捷金额列表。
// Convert top-up presets to the yuan-denominated response shape.
func toTopupPresetDTOs(presets []model.TopupPreset) []topupPresetDTO {
	out := make([]topupPresetDTO, 0, len(presets))
	for _, p := range presets {
		out = append(out, topupPresetDTO{
			Amount:     p.Amount,
			BonusQuota: quotaToYuan(p.BonusQuota),
		})
	}
	return out
}

// logStatisticDTO 按 day × model 聚合的日志统计输出：Quota 以「元」为单位。
// Log-statistic (grouped by day × model) response: Quota expressed in CNY yuan.
//
// 字段名保持 model.LogStatistic 原始的大写形式（该结构未声明 json tag），
// 避免在本次「只换单位」的变更里顺手改字段名而破坏既有调用方。
// Field names keep model.LogStatistic's original capitalized form (the struct declares
// no json tags) so this unit-only change does not break existing callers.
type logStatisticDTO struct {
	*model.LogStatistic
	Quota float64 `json:"Quota"`
}

// toLogStatisticDTO 把单行聚合统计转换为「元」口径的输出结构。
// Convert a single aggregated statistics row to the yuan-denominated response shape.
func toLogStatisticDTO(s *model.LogStatistic) *logStatisticDTO {
	if s == nil {
		return nil
	}
	return &logStatisticDTO{
		LogStatistic: s,
		Quota:        quotaToYuan(int64(s.Quota)),
	}
}

// toLogStatisticDTOs 批量转换聚合统计；入参为 nil 时返回空切片而非 nil。
// Convert aggregated statistics; returns an empty slice (not nil) when the input is empty.
func toLogStatisticDTOs(stats []*model.LogStatistic) []*logStatisticDTO {
	out := make([]*logStatisticDTO, 0, len(stats))
	for _, s := range stats {
		out = append(out, toLogStatisticDTO(s))
	}
	return out
}

// planInfoQuotaToYuan 将订单快照 plan_info 中的 bonus_quota 由微元改写为元。
// 订单快照 plan_info 是嵌套 JSON 字符串，无法用字段遮蔽处理，故单独改写：
//   - 仅当存在 bonus_quota 数值字段时改写（充值订单），其余结构原样返回；
//   - 用 map 承载以保留其余字段（amount / credit_amount 已经是元口径）；
//   - 解析失败时原样返回，不影响订单列表可用性。
//
// Rewrite bonus_quota inside the nested order snapshot JSON from micro-quota to yuan.
// Only top-up orders carry a numeric bonus_quota; anything else (or unparsable input)
// is returned unchanged.
func planInfoQuotaToYuan(planInfo string) string {
	if planInfo == "" {
		return planInfo
	}
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(planInfo), &fields); err != nil {
		return planInfo
	}
	changed := false
	// 充值订单：bonus_quota 微元 → 元
	if raw, ok := fields["bonus_quota"]; ok {
		if micro, ok := raw.(float64); ok {
			fields["bonus_quota"] = quotaToYuan(int64(micro))
			changed = true
		}
	}
	// 套餐订单：快照内的 virtual_amount 微元 → 元
	if raw, ok := fields["virtual_amount"]; ok {
		if micro, ok := raw.(float64); ok {
			fields["virtual_amount"] = quotaToYuan(int64(micro))
			changed = true
		}
	}
	if !changed {
		return planInfo
	}
	b, err := json.Marshal(fields)
	if err != nil {
		return planInfo
	}
	return string(b)
}

// orderDTO 订单输出：amount 本就是「元」；plan_info 快照内的 bonus_quota 改写为「元」。
// Order response: amount is already CNY; bonus_quota inside the plan_info snapshot is
// rewritten from micro-quota to CNY.
type orderDTO struct {
	*model.Order
	PlanInfo string `json:"plan_info"`
}

// toOrderDTO 把单个订单转换为「元」口径的输出结构。
// Convert a single order to the yuan-denominated response shape.
func toOrderDTO(o *model.Order) *orderDTO {
	if o == nil {
		return nil
	}
	return &orderDTO{
		Order:    o,
		PlanInfo: planInfoQuotaToYuan(o.PlanInfo),
	}
}

// toOrderDTOs 批量转换订单列表；入参为 nil 时返回空切片而非 nil。
// Convert an order list; returns an empty slice (not nil) when the input is empty.
func toOrderDTOs(orders []*model.Order) []*orderDTO {
	out := make([]*orderDTO, 0, len(orders))
	for _, o := range orders {
		out = append(out, toOrderDTO(o))
	}
	return out
}

// quotaStatDTO 日志统计输出：quota / normal_quota / subscription_quota 以「元」为单位。
// Log-stat response: quota / normal_quota / subscription_quota expressed in CNY yuan.
type quotaStatDTO struct {
	Quota             float64 `json:"quota"`
	NormalQuota       float64 `json:"normal_quota"`
	SubscriptionQuota float64 `json:"subscription_quota"`
}

// toQuotaStatDTO 把日志统计换算为「元」口径。
// Convert log statistics to the yuan-denominated shape.
func toQuotaStatDTO(quota, normalQuota, subscriptionQuota int64) quotaStatDTO {
	return quotaStatDTO{
		Quota:             quotaToYuan(quota),
		NormalQuota:       quotaToYuan(normalQuota),
		SubscriptionQuota: quotaToYuan(subscriptionQuota),
	}
}

// ---------------------------------------------------------------------------
// 输入 DTO（元 -> 微元）
//
// 写入方向同样在 controller 层收口：请求体先按「元」解析，再显式乘回微元后交给 model。
// 由于内嵌结构的同名字段在 Unmarshal 时不会被赋值（外层字段优先接收，其值保持零值），
// 必须显式调用 toModel/toToken/toRedemption 把元写回对应的微元字段。
//
// 这样做的另一个必要性：若直接解码到 model 的 int64 字段，客户端传小数元
// （如 10.5）会直接报 "cannot unmarshal number 10.5 into ... of type int64"。
//
// Input DTOs (yuan -> micro-quota).
// Requests are parsed in yuan and explicitly converted back before reaching model.
// Decoding straight into the model's int64 field would reject fractional yuan (e.g. 10.5).
// ---------------------------------------------------------------------------

// tokenRemainQuotaToMicro 将令牌剩余额度的「元」输入换算为微元。
// -1 是「无限额度」哨兵值（与 UnlimitedQuota 字段配套），原样透传不做换算，
// 否则会被换算成 -1_000_000 微元，与哨兵语义冲突。
//
// Convert token remain_quota from yuan to micro-quota. -1 is the "unlimited" sentinel
// paired with UnlimitedQuota and is passed through untouched.
func tokenRemainQuotaToMicro(y float64) int64 {
	if y == -1 {
		return -1
	}
	return yuanToQuota(y)
}

// updateUserRequest PUT /api/user/ 的请求体：quota 以「元」传入。
//
// Quota 用指针接收，以区分「未传 quota（不改额度）」与「显式清零 quota=0」：
// 零值在 GORM 的 Updates(struct) 中会被跳过，只有能分辨二者才能把零值补写落库。
// PUT /api/user/ request body: quota is expressed in CNY yuan.
// The pointer keeps "quota omitted" apart from an explicit zero.
type updateUserRequest struct {
	model.User
	Quota *float64 `json:"quota"`
}

// toModel 转成 model.User，并把元额度换算回微元；未传 quota 时保持 0，
// 由调用方根据 req.Quota 是否为 nil 决定是否补写。
// Convert to model.User with the yuan-denominated quota written back as micro-quota.
func (r *updateUserRequest) toModel() model.User {
	u := r.User
	if r.Quota != nil {
		u.Quota = yuanToQuota(*r.Quota)
	}
	return u
}

// tokenWriteRequest 令牌新增/修改请求体：remain_quota 以「元」传入。
// Token create/update request body: remain_quota is expressed in CNY yuan.
type tokenWriteRequest struct {
	model.Token
	RemainQuota float64 `json:"remain_quota"`
}

// toToken 转成 model.Token，并把元额度换算回微元。
// Convert to model.Token with the yuan-denominated quota written back as micro-quota.
func (r *tokenWriteRequest) toToken() model.Token {
	t := r.Token
	t.RemainQuota = tokenRemainQuotaToMicro(r.RemainQuota)
	return t
}

// redemptionWriteRequest 兑换码新增/修改请求体：quota 以「元」传入。
// Redemption create/update request body: quota is expressed in CNY yuan.
type redemptionWriteRequest struct {
	model.Redemption
	Quota float64 `json:"quota"`
}

// toRedemption 转成 model.Redemption，并把元额度换算回微元。
// Convert to model.Redemption with the yuan-denominated quota written back as micro-quota.
func (r *redemptionWriteRequest) toRedemption() model.Redemption {
	rd := r.Redemption
	rd.Quota = yuanToQuota(r.Quota)
	return rd
}

// ---------------------------------------------------------------------------
// 套餐（plan）DTO
//
// plan.virtual_amount 是「套餐虚拟余额」（存储层微元），API 边界一律以「元」呈现。
// 与 user/token/redemption 同一套遮蔽 + 显式回写模式。
// Plan DTOs: plan.virtual_amount is stored in micro-quota but exposed in CNY yuan.
// ---------------------------------------------------------------------------

// planDTO 套餐输出：virtual_amount 以「元」为单位。
type planDTO struct {
	*model.Plan
	VirtualAmount float64 `json:"virtual_amount"`
}

// toPlanDTO 把单个套餐转换为「元」口径的输出结构。
// Convert a single plan to the yuan-denominated response shape.
func toPlanDTO(p *model.Plan) *planDTO {
	if p == nil {
		return nil
	}
	return &planDTO{
		Plan:          p,
		VirtualAmount: quotaToYuan(p.VirtualAmount),
	}
}

// toPlanDTOs 批量转换套餐列表；入参为 nil 时返回空切片而非 nil。
// Convert a plan list; returns an empty slice (not nil) when the input is empty.
func toPlanDTOs(plans []*model.Plan) []*planDTO {
	out := make([]*planDTO, 0, len(plans))
	for _, p := range plans {
		out = append(out, toPlanDTO(p))
	}
	return out
}

// userPlanDTO 订阅输出：virtual_amount / used_amount 以「元」为单位，内嵌 plan 同样元化。
type userPlanDTO struct {
	*model.UserPlan
	VirtualAmount float64  `json:"virtual_amount"`
	UsedAmount    float64  `json:"used_amount"`
	Plan          *planDTO `json:"plan,omitempty"`
}

// toUserPlanDTO 把单条订阅转换为「元」口径的输出结构。
// Convert a single subscription to the yuan-denominated response shape.
func toUserPlanDTO(up *model.UserPlan) *userPlanDTO {
	if up == nil {
		return nil
	}
	return &userPlanDTO{
		UserPlan:      up,
		VirtualAmount: quotaToYuan(up.VirtualAmount),
		UsedAmount:    quotaToYuan(up.UsedAmount),
		Plan:          toPlanDTO(up.Plan),
	}
}

// toUserPlanDTOs 批量转换订阅列表；入参为 nil 时返回空切片而非 nil。
// Convert a subscription list; returns an empty slice (not nil) when the input is empty.
func toUserPlanDTOs(ups []*model.UserPlan) []*userPlanDTO {
	out := make([]*userPlanDTO, 0, len(ups))
	for _, up := range ups {
		out = append(out, toUserPlanDTO(up))
	}
	return out
}

// planWriteRequest 套餐新增/修改请求体：virtual_amount 以「元」传入。
type planWriteRequest struct {
	model.Plan
	VirtualAmount float64 `json:"virtual_amount"`
}

// toPlan 转成 model.Plan，并把元额度换算回微元。
// Convert to model.Plan with the yuan-denominated virtual amount written back as micro-quota.
func (r *planWriteRequest) toPlan() model.Plan {
	p := r.Plan
	p.VirtualAmount = yuanToQuota(r.VirtualAmount)
	return p
}
