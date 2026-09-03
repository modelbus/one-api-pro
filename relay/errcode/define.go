package errcode

type ErrorCodeMapping struct {
	ErrorType string
	ErrorCode string
}

var StatusCodeMappings = map[int]ErrorCodeMapping{
	400: {ErrorType: "invalid_request_error", ErrorCode: "bad_request"},
	401: {ErrorType: "authentication_error", ErrorCode: "invalid_api_key"},
	402: {ErrorType: "insufficient_quota", ErrorCode: "insufficient_quota"},
	403: {ErrorType: "permission_error", ErrorCode: "forbidden"},
	404: {ErrorType: "invalid_request_error", ErrorCode: "model_not_found"},
	422: {ErrorType: "invalid_request_error", ErrorCode: "validation_error"},
	429: {ErrorType: "rate_limit_error", ErrorCode: "rate_limit_exceeded"},
	500: {ErrorType: "server_error", ErrorCode: "internal_error"},
	502: {ErrorType: "server_error", ErrorCode: "bad_gateway"},
	503: {ErrorType: "server_error", ErrorCode: "service_unavailable"},
	529: {ErrorType: "rate_limit_error", ErrorCode: "overloaded"},
}

const (
	ActionPassthrough = "passthrough"
	ActionRetry      = "retry"
	ActionDisable    = "disable"
	ActionCooldown   = "cooldown"
	ActionReturn     = "return"

	CategoryPassthrough = "passthrough"
	CategoryDisable     = "disable"
	CategoryCooldown    = "cooldown"
	CategoryRetry       = "retry"
)

type StatusCodeConfig struct {
	Action          string `json:"action"`
	CooldownSeconds int    `json:"cooldown_seconds,omitempty"`
	Category        string `json:"category,omitempty"`
}

var DefaultStatusCodeConfigs = map[int]StatusCodeConfig{
	400: {Action: ActionPassthrough, Category: CategoryPassthrough},
	401: {Action: ActionDisable, Category: CategoryDisable},
	402: {Action: ActionDisable, Category: CategoryDisable},
	403: {Action: ActionDisable, Category: CategoryDisable},
	404: {Action: ActionPassthrough, Category: CategoryPassthrough},
	422: {Action: ActionPassthrough, Category: CategoryPassthrough},
	429: {Action: ActionCooldown, Category: CategoryCooldown},
	500: {Action: ActionRetry, Category: CategoryRetry},
	502: {Action: ActionRetry, Category: CategoryRetry},
	503: {Action: ActionRetry, Category: CategoryRetry},
	529: {Action: ActionCooldown, Category: CategoryCooldown},
}

func GetMapping(statusCode int) (ErrorCodeMapping, bool) {
	m, ok := StatusCodeMappings[statusCode]
	return m, ok
}

func GetDefaultConfig(statusCode int) StatusCodeConfig {
	cfg, ok := DefaultStatusCodeConfigs[statusCode]
	if !ok {
		if statusCode/100 == 5 {
			return StatusCodeConfig{Action: ActionRetry, Category: CategoryRetry}
		}
		if statusCode/100 == 4 {
			return StatusCodeConfig{Action: ActionPassthrough, Category: CategoryPassthrough}
		}
		return StatusCodeConfig{Action: ActionRetry, Category: CategoryRetry}
	}
	return cfg
}

// ShouldRetry reports whether the given error action should be followed by a
// retry on another channel. "disable" is retryable too: a channel that errors
// with 401/402/403 (invalid key / insufficient quota / forbidden) should be
// disabled AND the request should transparently fail over to another channel
// (or, when all same-model channels are exhausted, a fallback channel) instead
// of surfacing the raw error to the client.
func ShouldRetry(action string) bool {
	return action == ActionRetry || action == ActionCooldown || action == ActionDisable
}