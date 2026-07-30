package ratio

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/Leon-PanPan/one-api-pro/model"
)

func TestGetModelPriceWithSuffixHandling(t *testing.T) {
	Convey("GetModelPrice suffix handling", t, func() {
		cache := model.GetModelPriceCacheForTest()
		cache.Lock()
		model.SetModelPriceMapForTest(map[string]*model.ModelPrice{
			"qwen-max": {
				ModelName:   "qwen-max",
				InputPrice:  2.4,
				OutputPrice: 2.4,
				BillingType: model.BillingTypeToken,
				Enabled:     true,
			},
			"command-r": {
				ModelName:   "command-r",
				InputPrice:  0.5,
				OutputPrice: 1.5,
				BillingType: model.BillingTypeToken,
				Enabled:     true,
			},
			"gpt-4o": {
				ModelName:   "gpt-4o",
				InputPrice:  2.5,
				OutputPrice: 10.0,
				CachedPrice: 1.25,
				BillingType: model.BillingTypeToken,
				Enabled:     true,
			},
		})
		cache.Unlock()
		defer func() {
			cache.Lock()
			model.SetModelPriceMapForTest(nil)
			cache.Unlock()
		}()

		Convey("qwen-xxx-internet suffix stripped", func() {
			result, err := GetModelPrice("qwen-max-internet")
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(result.InputPrice, ShouldEqual, 2.4)
			So(result.Found, ShouldBeTrue)
		})

		Convey("command-xxx-internet suffix stripped", func() {
			result, err := GetModelPrice("command-r-internet")
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(result.InputPrice, ShouldEqual, 0.5)
		})

		Convey("normal model name works directly", func() {
			result, err := GetModelPrice("gpt-4o")
			So(err, ShouldBeNil)
			So(result.InputPrice, ShouldEqual, 2.5)
			So(result.OutputPrice, ShouldEqual, 10.0)
			So(result.CachedPrice, ShouldEqual, 1.25)
		})

		Convey("unknown model returns error", func() {
			result, err := GetModelPrice("nonexistent-model")
			So(err, ShouldNotBeNil)
			So(result, ShouldBeNil)
		})

		Convey("pattern fallback for versioned model", func() {
			result, err := GetModelPrice("gpt-4o-2024-08-06")
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(result.InputPrice, ShouldEqual, 2.5)
		})
	})
}

func TestModelPriceExists(t *testing.T) {
	Convey("ModelPriceExists", t, func() {
		cache := model.GetModelPriceCacheForTest()
		cache.Lock()
		model.SetModelPriceMapForTest(map[string]*model.ModelPrice{
			"gpt-4o": {
				ModelName:   "gpt-4o",
				InputPrice:  2.5,
				OutputPrice: 10.0,
				BillingType: model.BillingTypeToken,
				Enabled:     true,
			},
		})
		cache.Unlock()
		defer func() {
			cache.Lock()
			model.SetModelPriceMapForTest(nil)
			cache.Unlock()
		}()

		Convey("existing model returns true", func() {
			So(ModelPriceExists("gpt-4o"), ShouldBeTrue)
		})

		Convey("non-existing model returns false", func() {
			So(ModelPriceExists("unknown-model"), ShouldBeFalse)
		})
	})
}

func TestGetGroupDiscountDelegation(t *testing.T) {
	Convey("GetGroupDiscount delegates to model package", t, func() {
		cache := model.GetGroupPriceCacheForTest()
		cache.Lock()
		model.SetGroupPriceMapForTest(map[string]map[string]float64{
			"vip":     {"": 0.8, "gpt-4o": 0.5},
			"default": {"": 1.0},
		})
		cache.Unlock()
		defer func() {
			cache.Lock()
			model.SetGroupPriceMapForTest(nil)
			cache.Unlock()
		}()

		So(GetGroupDiscount("vip", "gpt-4o"), ShouldEqual, 0.5)
		So(GetGroupDiscount("vip", "other"), ShouldEqual, 0.8)
		So(GetGroupDiscount("default", ""), ShouldEqual, 1.0)
		So(GetGroupDiscount("unknown", ""), ShouldEqual, 1.0)
	})
}