package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/modelbus/one-api-pro/common/i18n"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

// validationFieldLabels 把结构体字段映射到「用户可读字段名」的文案键。
// 键使用 StructNamespace（形如 "User.Password"），避免不同结构体里的同名字段串味。
var validationFieldLabels = map[string]string{
	"User.Username":    "validation_field_username",
	"User.Password":    "validation_field_password",
	"User.DisplayName": "validation_field_display_name",
	"User.Email":       "validation_field_email",
}

// validationTagRules 把 validator 的 tag 映射到文案模板键；模板里的 %s 由 tag 参数（如 min=8 的 8）填充。
var validationTagRules = map[string]string{
	"required": "validation_rule_required",
	"min":      "validation_rule_min",
	"max":      "validation_rule_max",
	"len":      "validation_rule_len",
	"email":    "validation_rule_email",
	"gt":       "validation_rule_gt",
	"gte":      "validation_rule_gte",
	"lt":       "validation_rule_lt",
	"lte":      "validation_rule_lte",
	"oneof":    "validation_rule_oneof",
}

// ValidationMessage 把 validator 的字段校验错误转换成用户可读的提示文案。
//
// 例如 Password 的 min 校验失败时返回「密码长度不能少于 8 位」，而不是原始英文
// "Key: 'User.Password' Error:Field validation for 'Password' failed on the 'min' tag"。
// 无法识别的字段/规则统一回落为通用的「无效的输入」提示，绝不把英文原始错误透给前端。
func ValidationMessage(c *gin.Context, err error) string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) || len(validationErrs) == 0 {
		return i18n.Translate(c, "invalid_input")
	}
	messages := make([]string, 0, len(validationErrs))
	for _, fieldErr := range validationErrs {
		if message := validationFieldMessage(c, fieldErr); message != "" {
			messages = append(messages, message)
		}
	}
	if len(messages) == 0 {
		return i18n.Translate(c, "invalid_input")
	}
	separator := "; "
	if strings.HasPrefix(i18n.GetLang(c), "zh") {
		separator = "；"
	}
	return strings.Join(messages, separator)
}

// validationFieldMessage 生成单个字段错误的提示，识别不了时返回空串由调用方回落。
func validationFieldMessage(c *gin.Context, fieldErr validator.FieldError) string {
	ruleKey, ok := validationTagRules[fieldErr.Tag()]
	if !ok {
		return ""
	}
	rule := translateOrEmpty(c, ruleKey)
	if rule == "" {
		return ""
	}
	if strings.Contains(rule, "%s") {
		rule = fmt.Sprintf(rule, fieldErr.Param())
	}
	labelKey, ok := validationFieldLabels[fieldErr.StructNamespace()]
	if !ok {
		labelKey = "validation_field_generic"
	}
	return translateOrEmpty(c, labelKey) + rule
}

// translateOrEmpty 在文案缺失时返回空串，避免把裸的 i18n key 输出给前端。
func translateOrEmpty(c *gin.Context, key string) string {
	message := i18n.Translate(c, key)
	if message == key {
		return ""
	}
	return message
}
