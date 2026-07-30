package errcode

import "testing"

func TestGetMapping(t *testing.T) {
	tests := []struct {
		statusCode int
		wantType   string
		wantCode   string
		wantOk     bool
	}{
		{400, "invalid_request_error", "bad_request", true},
		{401, "authentication_error", "invalid_api_key", true},
		{402, "insufficient_quota", "insufficient_quota", true},
		{403, "permission_error", "forbidden", true},
		{429, "rate_limit_error", "rate_limit_exceeded", true},
		{500, "server_error", "internal_error", true},
		{529, "rate_limit_error", "overloaded", true},
		{418, "", "", false},
	}
	for _, tt := range tests {
		m, ok := GetMapping(tt.statusCode)
		if ok != tt.wantOk {
			t.Errorf("GetMapping(%d) ok = %v, want %v", tt.statusCode, ok, tt.wantOk)
		}
		if ok && (m.ErrorType != tt.wantType || m.ErrorCode != tt.wantCode) {
			t.Errorf("GetMapping(%d) = {type:%s, code:%s}, want {type:%s, code:%s}",
				tt.statusCode, m.ErrorType, m.ErrorCode, tt.wantType, tt.wantCode)
		}
	}
}

func TestGetDefaultConfig(t *testing.T) {
	tests := []struct {
		statusCode int
		wantAction string
	}{
		{400, ActionPassthrough},
		{401, ActionDisable},
		{402, ActionDisable},
		{429, ActionCooldown},
		{500, ActionRetry},
		{502, ActionRetry},
		{529, ActionCooldown},
	}
	for _, tt := range tests {
		cfg := GetDefaultConfig(tt.statusCode)
		if cfg.Action != tt.wantAction {
			t.Errorf("GetDefaultConfig(%d) action = %v, want %v", tt.statusCode, cfg.Action, tt.wantAction)
		}
	}
}

func TestShouldRetry(t *testing.T) {
	if !ShouldRetry(ActionRetry) {
		t.Error("ShouldRetry(ActionRetry) = false, want true")
	}
	if !ShouldRetry(ActionCooldown) {
		t.Error("ShouldRetry(ActionCooldown) = false, want true")
	}
	if ShouldRetry(ActionPassthrough) {
		t.Error("ShouldRetry(ActionPassthrough) = true, want false")
	}
	if ShouldRetry(ActionDisable) {
		t.Error("ShouldRetry(ActionDisable) = true, want false")
	}
}

func TestGetDefaultConfigUnknownCodes(t *testing.T) {
	cfg5xx := GetDefaultConfig(599)
	if cfg5xx.Action != ActionRetry {
		t.Errorf("GetDefaultConfig(599) action = %v, want %v", cfg5xx.Action, ActionRetry)
	}
	cfg4xx := GetDefaultConfig(418)
	if cfg4xx.Action != ActionPassthrough {
		t.Errorf("GetDefaultConfig(418) action = %v, want %v", cfg4xx.Action, ActionPassthrough)
	}
}