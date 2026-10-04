package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/common/i18n"
	"github.com/modelbus/one-api-pro/model"

	"github.com/gin-gonic/gin"
)

// yuanQuotaOptionKeys 是在 /api/option/ 上以「元」口径对外暴露的额度类 key。
//
// 存储（system_settings.value）与 config 内存态仍是微元；仅在 API 边界换算：
// GET 除 1e6、PUT 乘 1e6（见 quota_dto.go）。
// ChannelDisableThreshold 不在其中——它是上游账户余额阈值（USD），与 CNY 无关。
//
// Quota option keys exposed in CNY yuan at the /api/option/ boundary.
// Storage and the in-memory config stay in micro-quota; conversion happens only here.
var yuanQuotaOptionKeys = map[string]bool{
	"QuotaForNewUser":      true,
	"QuotaForInviter":      true,
	"QuotaForInvitee":      true,
	"QuotaRemindThreshold": true,
	"PreConsumedQuota":     true,
}

func GetOptions(c *gin.Context) {
	var options []*model.Option
	config.OptionMapRWMutex.Lock()
	for k, v := range config.OptionMap {
		if strings.HasSuffix(k, "Token") || strings.HasSuffix(k, "Secret") {
			continue
		}
		value := helper.Interface2String(v)
		// 额度类 key 在 API 边界由微元换算为「元」（见 quota_dto.go）。
		if yuanQuotaOptionKeys[k] {
			micro, perr := strconv.ParseInt(value, 10, 64)
			if perr == nil {
				value = strconv.FormatFloat(quotaToYuan(micro), 'f', -1, 64)
			}
		}
		options = append(options, &model.Option{
			Key:   k,
			Value: value,
		})
	}
	config.OptionMapRWMutex.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    options,
	})
	return
}

func UpdateOption(c *gin.Context) {
	var option model.Option
	err := json.NewDecoder(c.Request.Body).Decode(&option)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": i18n.Translate(c, "invalid_parameter"),
		})
		return
	}
	switch option.Key {
	case "Theme":
		if !config.ValidThemes[option.Value] {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无效的主题",
			})
			return
		}
	case "QuotaPerUnit":
		// QuotaPerUnit 已改为系统常量（1 元 = 1_000_000 额度），语义同微信支付的「分」，
		// 仅参与内部分计费/充值换算，不对外展示、不入库、不可修改（v0.0.24 2026-10-03）。
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "QuotaPerUnit 为系统常量，不可修改",
		})
		return
	case "GitHubOAuthEnabled":
		if option.Value == "true" && config.GitHubClientId == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 GitHub OAuth，请先填入 GitHub Client Id 以及 GitHub Client Secret！",
			})
			return
		}
	case "EmailDomainRestrictionEnabled":
		if option.Value == "true" && len(config.EmailDomainWhitelist) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用邮箱域名限制，请先填入限制的邮箱域名！",
			})
			return
		}
	case "WeChatAuthEnabled":
		if option.Value == "true" && config.WeChatServerAddress == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用微信登录，请先填入微信登录相关配置信息！",
			})
			return
		}
	case "TurnstileCheckEnabled":
		if option.Value == "true" && config.TurnstileSiteKey == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无法启用 Turnstile 校验，请先填入 Turnstile 校验相关配置信息！",
			})
			return
		}
	}
	// 额度类 key 以「元」传入，换算回微元后再落库（见 quota_dto.go）。
	if yuanQuotaOptionKeys[option.Key] {
		yuan, perr := strconv.ParseFloat(option.Value, 64)
		if perr != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": i18n.Translate(c, "invalid_parameter"),
			})
			return
		}
		option.Value = strconv.FormatInt(yuanToQuota(yuan), 10)
	}
	err = model.UpdateOption(option.Key, option.Value)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}
