package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/model"
)

// paymentConfigPayload is the body shape for PUT /api/setting/payment/{wechat|alipay|bank}.
// Files are sent as multipart/form-data; non-file fields are JSON in a
// `config` form value. Cert files are saved to data/payment/<method>/.
type paymentConfigPayload struct {
	Enabled bool                   `json:"enabled"`
	Config  map[string]interface{} `json:"config"`
}

// GetPaymentSettings returns the full payment settings bundle for the
// settings UI.
func GetPaymentSettings(c *gin.Context) {
	keys := []string{
		model.SystemSettingKeyWechatEnabled, model.SystemSettingKeyWechatConfig,
		model.SystemSettingKeyAlipayEnabled, model.SystemSettingKeyAlipayConfig,
		model.SystemSettingKeyBankEnabled, model.SystemSettingKeyBankConfig,
	}
	out := gin.H{}
	for _, k := range keys {
		s, err := model.GetSystemSetting(k)
		if err != nil {
			out[shortKey(k)] = gin.H{"enabled": false, "config": gin.H{}}
			continue
		}
		out[shortKey(k)] = gin.H{
			"enabled":     boolFromJSON(s.Value, "enabled"),
			"config":      jsonRaw(s.Value),
			"description": s.Description,
			"updated_at":  s.UpdatedAt,
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": out})
}

// PutPaymentMethod handles PUT /api/setting/payment/:method.
// :method is one of wechat, alipay, bank.
// Accepts multipart/form-data with two parts:
//   - "config": JSON string with {"enabled": bool, "config": {...}}
//   - "cert_file" / "key_file" / "private_key_file" / "public_key_file":
//     optional file uploads, written under data/payment/<method>/<basename>
// The path of the saved file replaces any "xxx_file" key in the JSON
// config so the file is associated with the settings row.
func PutPaymentMethod(c *gin.Context) {
	method := c.Param("method")
	if !isPaymentMethod(method) {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "未知支付方式: " + method})
		return
	}
	payload, err := readPaymentPayload(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "解析配置失败: " + err.Error()})
		return
	}

	// Save uploaded cert/key files. Path stored back into the config map.
	configMap := payload.Config
	if configMap == nil {
		configMap = map[string]interface{}{}
	}
	for _, field := range fileFieldsFor(method) {
		f, err := c.FormFile(field)
		if err != nil {
			continue // no upload for this field
		}
		dir := filepath.Join("data", "payment", method)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "创建目录失败: " + err.Error()})
			return
		}
		dst := filepath.Join(dir, filepath.Base(f.Filename))
		if err := c.SaveUploadedFile(f, dst); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "保存文件失败: " + err.Error()})
			return
		}
		configMap[field] = dst
	}

	// Persist enabled flag + config object back to system_settings.
	cfg := map[string]interface{}{"enabled": payload.Enabled}
	for k, v := range configMap {
		cfg[k] = v
	}
	cfgJSON, err := jsonMarshalString(cfg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "序列化配置失败: " + err.Error()})
		return
	}
	enabledJSON, _ := jsonMarshalString(map[string]interface{}{"enabled": payload.Enabled})
	switch method {
	case "wechat":
		if err := model.UpsertSystemSetting(model.SystemSettingKeyWechatEnabled, enabledJSON, model.SystemSettingCategoryPayment, "微信支付开关"); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := model.UpsertSystemSetting(model.SystemSettingKeyWechatConfig, cfgJSON, model.SystemSettingCategoryPayment, "微信支付参数"); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
	case "alipay":
		if err := model.UpsertSystemSetting(model.SystemSettingKeyAlipayEnabled, enabledJSON, model.SystemSettingCategoryPayment, "支付宝开关"); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := model.UpsertSystemSetting(model.SystemSettingKeyAlipayConfig, cfgJSON, model.SystemSettingCategoryPayment, "支付宝参数"); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
	case "bank":
		if err := model.UpsertSystemSetting(model.SystemSettingKeyBankEnabled, enabledJSON, model.SystemSettingCategoryPayment, "银行转账开关"); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := model.UpsertSystemSetting(model.SystemSettingKeyBankConfig, cfgJSON, model.SystemSettingCategoryPayment, "银行账户信息"); err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "已保存"})
}

// GetPlanSettings returns the plan-operations config (upgrade mode + topup switch).
func GetPlanSettings(c *gin.Context) {
	mode := model.GetSystemSettingString(model.SystemSettingKeyPlanUpgradeMode)
	if mode == "" {
		mode = model.OrderUpgradeModePriceDiff
	}
	allowTopup := model.GetSystemSettingString(model.SystemSettingKeyPlanAllowTopup) == "true"
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"upgrade_mode": mode,
			"allow_topup":  allowTopup,
		},
	})
}

type putPlanSettingsReq struct {
	UpgradeMode string `json:"upgrade_mode"`
	AllowTopup  bool   `json:"allow_topup"`
}

// PutPlanSettings updates the plan-operations config.
func PutPlanSettings(c *gin.Context) {
	var req putPlanSettingsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	if req.UpgradeMode != model.OrderUpgradeModePriceDiff && req.UpgradeMode != model.OrderUpgradeModeStack {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "upgrade_mode 必须是 price_diff 或 stack"})
		return
	}
	if err := model.UpsertSystemSetting(model.SystemSettingKeyPlanUpgradeMode,
		req.UpgradeMode,
		model.SystemSettingCategoryPlan, "套餐升级模式"); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := model.UpsertSystemSetting(model.SystemSettingKeyPlanAllowTopup,
		fmt.Sprintf("%t", req.AllowTopup),
		model.SystemSettingCategoryPlan, "是否允许余额充值（仅占位）"); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "已保存"})
}

// helpers --------------------------------------------------------

func readPaymentPayload(c *gin.Context) (*paymentConfigPayload, error) {
	// Accept either JSON body or multipart form-data with a "config" field
	// (for compatibility with the file-upload flow).
	ct := c.GetHeader("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		raw := c.PostForm("config")
		if raw == "" {
			return nil, errors.New("multipart 缺少 config 字段")
		}
		return parsePaymentPayloadJSON(raw)
	}
	var p paymentConfigPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func parsePaymentPayloadJSON(raw string) (*paymentConfigPayload, error) {
	// Try the "double JSON" shape: {"enabled":..., "config":{...}}
	// then the simple shape: {"value":"price_diff"} (used by plan settings)
	// and fall back to that.
	var p paymentConfigPayload
	if err := jsonUnmarshalString(raw, &p); err == nil && (p.Enabled || p.Config != nil) {
		return &p, nil
	}
	var alt struct {
		Value interface{} `json:"value"`
	}
	if err := jsonUnmarshalString(raw, &alt); err == nil {
		// The plan setting uses this shape; payment never does. Pass through
		// as {enabled:false, config:{value:...}} so caller can detect.
		p.Config = map[string]interface{}{"value": alt.Value}
		return &p, nil
	}
	return nil, errors.New("config 不是合法 JSON")
}

func isPaymentMethod(m string) bool {
	switch m {
	case "wechat", "alipay", "bank":
		return true
	}
	return false
}

func fileFieldsFor(method string) []string {
	switch method {
	case "wechat":
		return []string{"cert_file", "key_file"}
	case "alipay":
		return []string{"private_key_file", "public_key_file"}
	}
	return nil
}

// shortKey maps "payment.wechat.enabled" -> "wechat_enabled"
func shortKey(k string) string {
	return strings.ReplaceAll(strings.TrimPrefix(k, "payment."), ".", "_")
}

func boolFromJSON(s, field string) bool {
	if s == "" {
		return false
	}
	var v map[string]interface{}
	if err := jsonUnmarshalString(s, &v); err != nil {
		return false
	}
	if b, ok := v[field].(bool); ok {
		return b
	}
	return false
}

func jsonRaw(s string) interface{} {
	if s == "" {
		return map[string]interface{}{}
	}
	var v interface{}
	if err := jsonUnmarshalString(s, &v); err != nil {
		return map[string]interface{}{}
	}
	return v
}
