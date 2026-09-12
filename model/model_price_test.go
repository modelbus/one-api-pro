package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBillingTypeConstants(t *testing.T) {
	Convey("Billing type constants", t, func() {
		So(BillingTypeToken, ShouldEqual, "token")
		So(BillingTypeRequest, ShouldEqual, "request")
		So(BillingTypePerRequest, ShouldEqual, "per_request")
	})
}

func TestModelPriceStruct(t *testing.T) {
	Convey("ModelPrice struct", t, func() {
		mp := ModelPrice{
			ModelName:      "gpt-4o",
			InputPrice:     2.5,
			OutputPrice:    10.0,
			CachedPrice:    1.25,
			PerRequestPrice: 0,
			BillingType:    BillingTypeToken,
			Enabled:        true,
		}
		So(mp.ModelName, ShouldEqual, "gpt-4o")
		So(mp.InputPrice, ShouldEqual, 2.5)
		So(mp.OutputPrice, ShouldEqual, 10.0)
		So(mp.CachedPrice, ShouldEqual, 1.25)
		So(mp.BillingType, ShouldEqual, BillingTypeToken)
		So(mp.Enabled, ShouldBeTrue)
	})

	Convey("ModelPrice per_request type", t, func() {
		mp := ModelPrice{
			ModelName:      "dall-e-3",
			InputPrice:     0,
			OutputPrice:    0,
			CachedPrice:    0,
			PerRequestPrice: 0.04,
			BillingType:    BillingTypePerRequest,
			Enabled:        true,
		}
		So(mp.BillingType, ShouldEqual, BillingTypePerRequest)
		So(mp.PerRequestPrice, ShouldEqual, 0.04)
	})
}

func TestGroupPriceStruct(t *testing.T) {
	Convey("GroupPrice struct", t, func() {
		gp := GroupPrice{
			GroupName: "vip",
			ModelName: "",
			Discount:  0.8,
		}
		So(gp.GroupName, ShouldEqual, "vip")
		So(gp.ModelName, ShouldEqual, "")
		So(gp.Discount, ShouldEqual, 0.8)
	})

	Convey("GroupPrice model-specific discount", t, func() {
		gp := GroupPrice{
			GroupName: "vip",
			ModelName: "gpt-4o",
			Discount:  0.5,
		}
		So(gp.ModelName, ShouldEqual, "gpt-4o")
		So(gp.Discount, ShouldEqual, 0.5)
	})
}

func TestGetGroupDiscount(t *testing.T) {
	Convey("GetGroupDiscount", t, func() {
		groupPriceCache.Lock()
		groupPriceMap = map[string]map[string]float64{
			"default": {"": 1.0},
			"vip":     {"": 0.8, "gpt-4o": 0.5},
		}
		groupPriceCache.Unlock()

		Convey("exact model match takes priority", func() {
			discount := GetGroupDiscount("vip", "gpt-4o")
			So(discount, ShouldEqual, 0.5)
		})

		Convey("default discount when no model match", func() {
			discount := GetGroupDiscount("vip", "claude-3")
			So(discount, ShouldEqual, 0.8)
		})

		Convey("returns 1.0 for unknown group", func() {
			discount := GetGroupDiscount("unknown", "")
			So(discount, ShouldEqual, 1.0)
		})

		Convey("returns 1.0 when cache is nil", func() {
			groupPriceCache.Lock()
			groupPriceMap = nil
			groupPriceCache.Unlock()
			discount := GetGroupDiscount("default", "")
			So(discount, ShouldEqual, 1.0)
		})
	})
}

func TestGetModelPrice(t *testing.T) {
	Convey("GetModelPrice", t, func() {
		modelPriceCache.Lock()
		modelPriceMap = map[string]*ModelPrice{
			"gpt-4o": {
				ModelName:   "gpt-4o",
				InputPrice:  2.5,
				OutputPrice: 10.0,
				CachedPrice: 1.25,
				BillingType: BillingTypeToken,
				Enabled:     true,
			},
			"dall-e-3": {
				ModelName:       "dall-e-3",
				PerRequestPrice: 0.04,
				BillingType:     BillingTypePerRequest,
				Enabled:         true,
			},
		}
		modelPriceCache.Unlock()

		Convey("returns price for known model", func() {
			mp, found := GetModelPrice("gpt-4o")
			So(found, ShouldBeTrue)
			So(mp.ModelName, ShouldEqual, "gpt-4o")
			So(mp.InputPrice, ShouldEqual, 2.5)
			So(mp.OutputPrice, ShouldEqual, 10.0)
		})

		Convey("returns not found for unknown model", func() {
			_, found := GetModelPrice("unknown-model")
			So(found, ShouldBeFalse)
		})

		Convey("per_request model type", func() {
			mp, found := GetModelPrice("dall-e-3")
			So(found, ShouldBeTrue)
			So(mp.BillingType, ShouldEqual, BillingTypePerRequest)
			So(mp.PerRequestPrice, ShouldEqual, 0.04)
		})

		Convey("returns not found when cache is nil", func() {
			modelPriceCache.Lock()
			modelPriceMap = nil
			modelPriceCache.Unlock()
			_, found := GetModelPrice("gpt-4o")
			So(found, ShouldBeFalse)
		})
	})
}

func TestFindModelPriceByPattern(t *testing.T) {
	Convey("FindModelPriceByPattern", t, func() {
		modelPriceCache.Lock()
		modelPriceMap = map[string]*ModelPrice{
			"gpt-4o": {
				ModelName:   "gpt-4o",
				InputPrice:  2.5,
				OutputPrice: 10.0,
				BillingType: BillingTypeToken,
				Enabled:     true,
			},
		}
		modelPriceCache.Unlock()

		Convey("exact match first", func() {
			mp := FindModelPriceByPattern("gpt-4o")
			So(mp, ShouldNotBeNil)
			So(mp.ModelName, ShouldEqual, "gpt-4o")
		})

		Convey("pattern match by substring", func() {
			mp := FindModelPriceByPattern("gpt-4o-2024-08-06")
			So(mp, ShouldNotBeNil)
			So(mp.ModelName, ShouldEqual, "gpt-4o")
		})

		Convey("no match returns nil", func() {
			mp := FindModelPriceByPattern("claude-3-opus")
			So(mp, ShouldBeNil)
		})

		Convey("nil cache returns nil", func() {
			modelPriceCache.Lock()
			modelPriceMap = nil
			modelPriceCache.Unlock()
			mp := FindModelPriceByPattern("gpt-4o")
			So(mp, ShouldBeNil)
		})
	})
}

func TestGetGroupNames(t *testing.T) {
	Convey("GetGroupNames", t, func() {
		groupPriceCache.Lock()
		groupPriceMap = map[string]map[string]float64{
			"default": {"": 1.0},
			"vip":     {"": 0.8},
			"svip":    {"": 0.6},
		}
		groupPriceCache.Unlock()

		names := GetGroupNames()
		So(len(names), ShouldEqual, 3)

		Convey("nil cache returns empty", func() {
			groupPriceCache.Lock()
			groupPriceMap = nil
			groupPriceCache.Unlock()
			names := GetGroupNames()
			So(len(names), ShouldEqual, 0)
		})
	})
}

func TestDefaultModelPrices(t *testing.T) {
	Convey("Default model prices", t, func() {
		So(len(defaultModelPrices), ShouldBeGreaterThan, 0)

		names := make(map[string]bool)
		for _, p := range defaultModelPrices {
			Convey(p.ModelName+" should have valid billing type", func() {
				So(p.BillingType, ShouldBeIn, BillingTypeToken, BillingTypePerRequest)
			})
			Convey(p.ModelName+" should have unique name", func() {
				So(names[p.ModelName], ShouldBeFalse)
				names[p.ModelName] = true
			})
		}
	})

	Convey("Default group prices", t, func() {
		So(len(defaultGroupPrices), ShouldBeGreaterThan, 0)

		for _, p := range defaultGroupPrices {
			So(p.GroupName, ShouldNotBeEmpty)
			So(p.Discount, ShouldBeGreaterThan, 0)
		}
	})
}

// setupPriceTestDB opens an in-memory sqlite + model_prices/group_prices tables
// and wires it as the package-level DB. Returns the *gorm.DB for convenience.
//
// 版本: v0.0.18
// 日期: 2026-09-13
func setupPriceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&ModelPrice{}, &GroupPrice{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// TestModelPriceInsert_StripsClientId 验证 Insert 不会使用 struct 上预置的
// 非零主键，从而避免「前端编辑后新增时携带旧 id 写入数据库」造成的冲突。
//
// 防御深度：即便 controller 未清零 Id，model 层也应保证主键自增。
// 由 controller.AddModelPrice 的 price.Id = 0 兜底，本测试聚焦 model.Insert() 行为。
//
// 版本: v0.0.18
// 日期: 2026-09-13
func TestModelPriceInsert_StripsClientId(t *testing.T) {
	Convey("ModelPrice.Insert ignores preset Id", t, func() {
		db := setupPriceTestDB(t)
		DB = db

		p := &ModelPrice{
			Id:          9999,
			ModelName:   "gpt-4o",
			InputPrice:  2.5,
			OutputPrice: 10.0,
			BillingType: BillingTypeToken,
			Enabled:     true,
		}
		So(p.Insert(), ShouldBeNil)

		Convey("DB assigns auto-incremented primary key, not the preset", func() {
			So(p.Id, ShouldNotEqual, 9999)
			So(p.Id, ShouldBeGreaterThan, 0)
		})

		Convey("CreatedAt/UpdatedAt filled", func() {
			So(p.CreatedAt, ShouldBeGreaterThan, 0)
			So(p.UpdatedAt, ShouldBeGreaterThan, 0)
		})

		Convey("row persisted with model_name preserved", func() {
			var got ModelPrice
			So(DB.First(&got, "model_name = ?", "gpt-4o").Error, ShouldBeNil)
			So(got.Id, ShouldEqual, p.Id)
			So(got.InputPrice, ShouldEqual, 2.5)
		})
	})
}

// TestGroupPriceInsert_StripsClientId 同上的 GroupPrice 版本。
//
// 版本: v0.0.18
// 日期: 2026-09-13
func TestGroupPriceInsert_StripsClientId(t *testing.T) {
	Convey("GroupPrice.Insert ignores preset Id", t, func() {
		db := setupPriceTestDB(t)
		DB = db

		p := &GroupPrice{
			Id:        9999,
			GroupName: "vip",
			ModelName: "gpt-4o",
			Discount:  0.5,
		}
		So(p.Insert(), ShouldBeNil)

		Convey("DB assigns auto-incremented primary key, not the preset", func() {
			So(p.Id, ShouldNotEqual, 9999)
			So(p.Id, ShouldBeGreaterThan, 0)
		})

		Convey("CreatedAt/UpdatedAt filled", func() {
			So(p.CreatedAt, ShouldBeGreaterThan, 0)
			So(p.UpdatedAt, ShouldBeGreaterThan, 0)
		})

		Convey("row persisted with group_name+model_name preserved", func() {
			var got GroupPrice
			So(DB.First(&got, "group_name = ? AND model_name = ?", "vip", "gpt-4o").Error, ShouldBeNil)
			So(got.Id, ShouldEqual, p.Id)
			So(got.Discount, ShouldEqual, 0.5)
		})
	})
}