package model

import "testing"

func TestChannelContainsGroup(t *testing.T) {
	ch := &Channel{Group: "default,vip"}

	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty group always allowed", "", true},
		{"single exact match default", "default", true},
		{"single exact match vip", "vip", true},
		{"group not in list", "premium", false},
		{"substring should not match", "def", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ch.ContainsGroup(tc.input); got != tc.want {
				t.Errorf("ContainsGroup(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestChannelContainsModel(t *testing.T) {
	ch := &Channel{Models: "gpt-4,gpt-3.5-turbo"}

	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty model always allowed", "", true},
		{"exact match gpt-4", "gpt-4", true},
		{"exact match gpt-3.5-turbo", "gpt-3.5-turbo", true},
		{"model not in list", "claude-3", false},
		{"substring should not match", "gpt", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ch.ContainsModel(tc.input); got != tc.want {
				t.Errorf("ContainsModel(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}