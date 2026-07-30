package ratio

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/Leon-PanPan/one-api-pro/common/config"
)

func TestCalculateTokenQuota(t *testing.T) {
	Convey("CalculateTokenQuota", t, func() {
		originalQuotaPerUnit := config.QuotaPerUnit
		config.QuotaPerUnit = 500000

		defer func() {
			config.QuotaPerUnit = originalQuotaPerUnit
		}()

		Convey("basic token billing", func() {
			quota := CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 0, 1.0)
			So(quota, ShouldBeGreaterThan, 0)
		})

		Convey("zero tokens should return 0", func() {
			quota := CalculateTokenQuota(2.5, 10, 1.25, 0, 0, 0, 1.0)
			So(quota, ShouldEqual, 0)
		})

		Convey("cached tokens reduce input tokens", func() {
			quotaNoCache := CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 0, 1.0)
			quotaWithCache := CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 500, 1.0)
			So(quotaWithCache, ShouldBeLessThan, quotaNoCache)
		})

		Convey("cached tokens at lower price", func() {
			quotaAllCached := CalculateTokenQuota(2.5, 10, 1.25, 1000, 0, 1000, 1.0)
			expected := int64(1.25 * 1000 / 1_000_000 * 500000)
			So(quotaAllCached, ShouldEqual, expected)
		})

		Convey("group discount applies", func() {
			quotaFull := CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 0, 1.0)
			quotaHalf := CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 0, 0.5)
			So(quotaHalf, ShouldEqual, quotaFull/2)
		})

		Convey("zero group discount defaults to 1.0", func() {
			quotaNormal := CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 0, 1.0)
			quotaZero := CalculateTokenQuota(2.5, 10, 1.25, 1000, 500, 0, 0)
			So(quotaZero, ShouldEqual, quotaNormal)
		})

		Convey("minimum quota is 1 for non-zero tokens", func() {
			quota := CalculateTokenQuota(0.001, 0.001, 0, 1, 1, 0, 1.0)
			So(quota, ShouldEqual, 1)
		})

		Convey("gpt-4o pricing example", func() {
			inputPrice := 2.5
			outputPrice := 10.0
			cachedPrice := 1.25
			quota := CalculateTokenQuota(inputPrice, outputPrice, cachedPrice, 1000000, 1000000, 0, 1.0)
			expected := int64((2.5*1000000+10.0*1000000)*1.0/1_000_000*500000)
			So(quota, ShouldEqual, expected)
		})

		Convey("gpt-4o with cache pricing example", func() {
			inputPrice := 2.5
			outputPrice := 10.0
			cachedPrice := 1.25
			quota := CalculateTokenQuota(inputPrice, outputPrice, cachedPrice, 1000000, 500000, 500000, 1.0)
			inputTokens := 1000000 - 500000
			expected := int64((2.5*float64(inputTokens)+10.0*500000+1.25*500000)*1.0/1_000_000*500000)
			So(quota, ShouldEqual, expected)
		})
	})
}

func TestCalculatePerRequestQuota(t *testing.T) {
	Convey("CalculatePerRequestQuota", t, func() {
		originalQuotaPerUnit := config.QuotaPerUnit
		config.QuotaPerUnit = 500000

		defer func() {
			config.QuotaPerUnit = originalQuotaPerUnit
		}()

		Convey("basic per request billing", func() {
			quota := CalculatePerRequestQuota(0.04, 1, 1, 1.0)
			So(quota, ShouldBeGreaterThan, 0)
		})

		Convey("size ratio applies", func() {
			quota1x := CalculatePerRequestQuota(0.04, 1, 1, 1.0)
			quota2x := CalculatePerRequestQuota(0.04, 2, 1, 1.0)
			So(quota2x, ShouldEqual, quota1x*2)
		})

		Convey("N count applies", func() {
			quota1 := CalculatePerRequestQuota(0.04, 1, 1, 1.0)
			quota3 := CalculatePerRequestQuota(0.04, 1, 3, 1.0)
			So(quota3, ShouldEqual, quota1*3)
		})

		Convey("group discount applies", func() {
			quota1 := CalculatePerRequestQuota(0.04, 1, 1, 1.0)
			quota08 := CalculatePerRequestQuota(0.04, 1, 1, 0.8)
			So(quota08, ShouldEqual, int64(float64(quota1)*0.8))
		})

		Convey("zero group discount defaults to 1.0", func() {
			quota1 := CalculatePerRequestQuota(0.04, 1, 1, 1.0)
			quota0 := CalculatePerRequestQuota(0.04, 1, 1, 0)
			So(quota0, ShouldEqual, quota1)
		})

		Convey("zero N defaults to 1", func() {
			quota0 := CalculatePerRequestQuota(0.04, 1, 0, 1.0)
			quota1 := CalculatePerRequestQuota(0.04, 1, 1, 1.0)
			So(quota0, ShouldEqual, quota1)
		})

		Convey("dall-e-3 pricing example", func() {
			quota := CalculatePerRequestQuota(0.04, 1, 1, 1.0)
			expected := int64(0.04 * 1 * 1 * 1.0 * 500000)
			So(quota, ShouldEqual, expected)
		})

		Convey("minimum quota is 1", func() {
			quota := CalculatePerRequestQuota(0.000001, 1, 1, 1.0)
			So(quota, ShouldEqual, 1)
		})
	})
}

func TestQuotaToUnit(t *testing.T) {
	Convey("QuotaToUnit", t, func() {
		originalQuotaPerUnit := config.QuotaPerUnit
		config.QuotaPerUnit = 500000

		defer func() {
			config.QuotaPerUnit = originalQuotaPerUnit
		}()

		Convey("converts quota to unit", func() {
			unit := QuotaToUnit(500000)
			So(unit, ShouldEqual, 1.0)
		})

		Convey("zero quota", func() {
			unit := QuotaToUnit(0)
			So(unit, ShouldEqual, 0)
		})
	})
}

func TestMillionConstant(t *testing.T) {
	Convey("Million constant", t, func() {
		So(Million, ShouldEqual, 1_000_000.0)
	})
}