package model

import "testing"

func TestCalculateUpgradePrice(t *testing.T) {
	now := int64(1700000000) // some fixed timestamp
	day := int64(86400)

	cases := []struct {
		name       string
		oldPrice   float64
		oldEndTime int64
		expired    bool
		newPrice   float64
		newSort    int
		oldSort    int
		want       float64
	}{
		{
			name: "expired old plan returns full new price",
			// old plan expired 1 day ago
			oldPrice: 9.9, oldEndTime: now - day, newPrice: 29.9, want: 29.9,
		},
		{
			name: "30 days remaining, new is pricier",
			oldPrice: 9.9, oldEndTime: now + 30*day, newPrice: 29.9, want: 20.0, // (29.9-9.9)/30 * 30 = 20
		},
		{
			name: "30 days remaining, new is cheaper (no refund)",
			oldPrice: 29.9, oldEndTime: now + 30*day, newPrice: 9.9, want: 0.0,
		},
		{
			name: "1 day remaining, small diff",
			oldPrice: 9.9, oldEndTime: now + day, newPrice: 19.9, want: 0.33, // (19.9-9.9)/30 = 0.33
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			oldPlan := &Plan{Id: 1, Name: "basic", Price: c.oldPrice, Sort: 1}
			newPlan := &Plan{Id: 2, Name: "pro", Price: c.newPrice, Sort: 5}
			got := CalculateUpgradePrice(oldPlan, c.oldEndTime, now, newPlan)
			// Allow 1 cent of rounding tolerance.
			diff := got - c.want
			if diff < -0.02 || diff > 0.02 {
				t.Errorf("got %v, want ~%v (delta %v)", got, c.want, diff)
			}
		})
	}
}

func TestIsValidPayMethod(t *testing.T) {
	valid := []string{"wechat", "alipay", "bank", "offline", "free", "WeChat", "WECHAT"}
	for _, m := range valid {
		if !IsValidPayMethod(m) {
			t.Errorf("expected %q to be valid", m)
		}
	}
	invalid := []string{"", "paypal", "stripe", "wechat-pay"}
	for _, m := range invalid {
		if IsValidPayMethod(m) {
			t.Errorf("expected %q to be invalid", m)
		}
	}
}
