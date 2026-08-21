package ratio

import (
	"fmt"
	"strings"

	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/model"
)

const Million = 1_000_000.0

type PriceResult struct {
	InputPrice      float64
	OutputPrice     float64
	CachedPrice     float64
	PerRequestPrice float64
	BillingType     string
	Discount        float64
	Found           bool
}

func GetModelPrice(modelName string, fallbackModelNames ...string) (*PriceResult, error) {
	tryModelName := modelName
	if strings.HasPrefix(tryModelName, "qwen-") && strings.HasSuffix(tryModelName, "-internet") {
		tryModelName = strings.TrimSuffix(tryModelName, "-internet")
	}
	if strings.HasPrefix(tryModelName, "command-") && strings.HasSuffix(tryModelName, "-internet") {
		tryModelName = strings.TrimSuffix(tryModelName, "-internet")
	}

	mp, found := model.GetModelPrice(tryModelName)
	if !found {
		mp = model.FindModelPriceByPattern(tryModelName)
	}
	if mp == nil {
		for _, fallback := range fallbackModelNames {
			fb := fallback
			if strings.HasPrefix(fb, "qwen-") && strings.HasSuffix(fb, "-internet") {
				fb = strings.TrimSuffix(fb, "-internet")
			}
			if strings.HasPrefix(fb, "command-") && strings.HasSuffix(fb, "-internet") {
				fb = strings.TrimSuffix(fb, "-internet")
			}
			mp, found = model.GetModelPrice(fb)
			if !found {
				mp = model.FindModelPriceByPattern(fb)
			}
			if mp != nil {
				break
			}
		}
	}
	if mp == nil {
		names := append([]string{modelName}, fallbackModelNames...)
		return nil, fmt.Errorf("model price not found for model(s): %v; please add a price entry for either the mapped model name or the original request model name", names)
	}
	return &PriceResult{
		InputPrice:      mp.InputPrice,
		OutputPrice:     mp.OutputPrice,
		CachedPrice:     mp.CachedPrice,
		PerRequestPrice: mp.PerRequestPrice,
		BillingType:     mp.BillingType,
		Found:           true,
	}, nil
}

func GetGroupDiscount(groupName string, modelName string, fallbackModelNames ...string) float64 {
	discount := model.GetGroupDiscount(groupName, modelName)
	if discount != 1.0 {
		return discount
	}
	for _, fallback := range fallbackModelNames {
		discount = model.GetGroupDiscount(groupName, fallback)
		if discount != 1.0 {
			return discount
		}
	}
	return 1.0
}

func CalculateTokenQuota(inputPrice float64, outputPrice float64, cachedPrice float64, promptTokens int, completionTokens int, cachedTokens int, groupDiscount float64) int64 {
	if groupDiscount == 0 {
		groupDiscount = 1.0
	}
	inputTokens := promptTokens - cachedTokens
	quota := (inputPrice*float64(inputTokens) + outputPrice*float64(completionTokens) + cachedPrice*float64(cachedTokens)) * groupDiscount / Million * config.QuotaPerUnit
	result := int64(quota)
	if result <= 0 && (promptTokens+completionTokens) > 0 {
		result = 1
	}
	return result
}

func CalculatePerRequestQuota(perRequestPrice float64, sizeRatio float64, n int, groupDiscount float64) int64 {
	if groupDiscount == 0 {
		groupDiscount = 1.0
	}
	if n <= 0 {
		n = 1
	}
	quota := perRequestPrice * sizeRatio * float64(n) * groupDiscount * config.QuotaPerUnit
	result := int64(quota)
	if result <= 0 {
		result = 1
	}
	return result
}

func QuotaToUnit(quota int64) float64 {
	return float64(quota) / config.QuotaPerUnit
}

func ModelPriceExists(modelName string) bool {
	_, err := GetModelPrice(modelName)
	return err == nil
}

// GetModelRatio remains for backward compatibility with billing.go PostConsumeQuota
func GetModelRatio(name string, channelType int) float64 {
	logger.SysError("GetModelRatio is deprecated, use GetModelPrice instead")
	return 1
}

func GetCompletionRatio(name string, channelType int) float64 {
	logger.SysError("GetCompletionRatio is deprecated, use GetModelPrice instead")
	return 1
}