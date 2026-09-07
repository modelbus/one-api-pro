package model

import (
	"testing"
)

// TestStringSlice_MarshalJSON 输出 JSON 数组；nil 输出 [] 而非 null。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func TestStringSlice_MarshalJSON(t *testing.T) {
	cases := []struct {
		in   StringSlice
		want string
	}{
		{nil, "[]"},
		{StringSlice{}, "[]"},
		{StringSlice{"A"}, `["A"]`},
		{StringSlice{"A", "B", "C"}, `["A","B","C"]`},
	}
	for _, c := range cases {
		got, err := c.in.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON(%v) error: %v", c.in, err)
		}
		if string(got) != c.want {
			t.Errorf("MarshalJSON(%v) = %s, want %s", c.in, got, c.want)
		}
	}
}

// TestStringSlice_UnmarshalJSON 覆盖 JSON 数组 / 旧字符串 / 换行文本 / null / 空。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func TestStringSlice_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want StringSlice
	}{
		{"empty array", `[]`, StringSlice{}},
		{"array", `["A","B","C"]`, StringSlice{"A", "B", "C"}},
		{"array with whitespace", `  ["A", "B"]  `, StringSlice{"A", "B"}},
		{"legacy newline text", `"A\nB\nC"`, StringSlice{"A", "B", "C"}},
		{"legacy single line", `"A"`, StringSlice{"A"}},
		{"legacy CRLF text", `"A\r\nB"`, StringSlice{"A", "B"}},
		{"null", `null`, nil},
		{"empty string", `""`, StringSlice{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got StringSlice
			if err := got.UnmarshalJSON([]byte(c.in)); err != nil {
				t.Fatalf("UnmarshalJSON(%s) error: %v", c.in, err)
			}
			if !equalSlice(got, c.want) {
				t.Errorf("UnmarshalJSON(%s) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

// TestStringSlice_UnmarshalJSON_Invalid 错误格式应返回 error。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func TestStringSlice_UnmarshalJSON_Invalid(t *testing.T) {
	var s StringSlice
	if err := s.UnmarshalJSON([]byte("not-a-json")); err == nil {
		t.Errorf("expected error for invalid input, got nil")
	}
}

// TestStringSlice_Value_Scan 验证 DB ↔ Go 的双向转换，涵盖空 / 数组 / 旧换行。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func TestStringSlice_Value_Scan(t *testing.T) {
	cases := []struct {
		name   string
		raw    string
		expect StringSlice
	}{
		{"empty", "", nil},
		{"json array", `["A","B"]`, StringSlice{"A", "B"}},
		{"legacy newline", "A\nB\nC", StringSlice{"A", "B", "C"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got StringSlice
			if err := got.Scan(tc.raw); err != nil {
				t.Fatalf("Scan(%q) error: %v", tc.raw, err)
			}
			if !equalSlice(got, tc.expect) {
				t.Errorf("Scan(%q) = %v, want %v", tc.raw, got, tc.expect)
			}
		})
	}
}

// TestStringSlice_Scan_Nil nil 入参不应 panic 且结果为 nil。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func TestStringSlice_Scan_Nil(t *testing.T) {
	var s StringSlice
	if err := s.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error: %v", err)
	}
	if s != nil {
		t.Errorf("Scan(nil) = %v, want nil", s)
	}
}

// TestStringSlice_Value_Nil nil 输出应为 nil driver.Value。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func TestStringSlice_Value_Nil(t *testing.T) {
	var s StringSlice
	v, err := s.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if v != nil {
		t.Errorf("Value() on nil = %v, want nil", v)
	}
}

// TestStringSlice_RoundTrip MarshalJSON → UnmarshalJSON 应保真。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func TestStringSlice_RoundTrip(t *testing.T) {
	original := StringSlice{"A", "B", "C"}
	data, err := original.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	var got StringSlice
	if err := got.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}
	if !equalSlice(got, original) {
		t.Errorf("round trip mismatch: got %v, want %v", got, original)
	}
}

// equalSlice 比较两个 StringSlice（nil 与空 slice 视为相等）。
//
// 版本: v0.0.12
// 日期: 2026-09-07
// 作者: opencode
func equalSlice(a, b StringSlice) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}