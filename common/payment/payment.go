// Package payment provides a thin abstraction over multiple Chinese
// payment channels (WeChat Native, Alipay Face-to-Face, etc.) and is
// consumed by the order module to obtain pre-payment URLs / QR codes
// and to verify async payment notifications.
//
// Each channel implements the Channel interface and is constructed via
// New(payMethod) by reading credentials from the system_settings table.
package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/modelbus/one-api-pro/model"
)

// PrePayResult is returned by Channel.PrePay.
type PrePayResult struct {
	PayURL   string `json:"pay_url"`   // URL the user should be redirected to, or QR-encoded string
	QRCode   string `json:"qr_code"`   // Same as PayURL for Native channels (caller decides how to render)
	ExpireAt int64  `json:"expire_at"` // Unix seconds; 0 means unknown
	TradeNo  string `json:"trade_no"`  // Provider's prepay id (separate from our order_no)
}

// Channel is implemented by every supported payment channel.
type Channel interface {
	// Name returns the stable identifier (matches OrderPayMethod*).
	Name() string
	// IsEnabled returns whether the channel is currently usable.
	// Implementations read from system_settings at call time so that
	// toggling the admin UI takes effect immediately.
	IsEnabled() (bool, error)
	// PrePay builds a pre-payment request. The returned result is opaque
	// to the caller; payment.go forwards it to the order API consumer.
	PrePay(orderNo string, amount float64, subject string) (*PrePayResult, error)
	// VerifyNotify verifies the signature of a payment-channel async
	// notification. The channel reads its own notify payload and returns
	// the parsed business fields on success.
	// For wechat: form-encoded body with sign verification.
	// For alipay: form-encoded body with RSA2 sign verification.
	VerifyNotify(payload []byte) (*NotifyResult, error)
}

// NotifyResult is the channel-agnostic shape of a successful async
// notification.
type NotifyResult struct {
	OutTradeNo string  // equals our order_no
	TradeNo   string  // payment provider's transaction id
	Amount     float64
	Paid       bool
}

var (
	registryMu sync.RWMutex
	registry   = map[string]Channel{}
)

// RegisterChannel makes a channel available to New().
func RegisterChannel(c Channel) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[c.Name()] = c
}

// New returns the registered channel for the given pay_method, or an
// error if no such channel is registered.
func New(payMethod string) (Channel, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	c, ok := registry[payMethod]
	if !ok {
		return nil, fmt.Errorf("未注册的支付方式: %s", payMethod)
	}
	return c, nil
}

// HasChannel reports whether a channel has been registered for the
// given pay_method.
func HasChannel(payMethod string) bool {
	registryMu.RLock()
	defer registryMu.RUnlock()
	_, ok := registry[payMethod]
	return ok
}

// PaymentMethod describes one payment channel's enabled status for
// user-facing display. Methods that are not registered are omitted.
type PaymentMethod struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Enabled bool   `json:"enabled"`
}

// PaymentMethodLabel returns a human-readable label for a pay_method.
func PaymentMethodLabel(name string) string {
	switch name {
	case model.OrderPayMethodWechat:
		return "微信支付"
	case model.OrderPayMethodAlipay:
		return "支付宝"
	case model.OrderPayMethodBank:
		return "银行转账"
	case model.OrderPayMethodOffline:
		return "线下支付"
	default:
		return name
	}
}

// AnyChannelEnabled returns true when at least one of the registered
// channels reports IsEnabled()==true. It is used to short-circuit the
// user-facing purchase flow when the admin has not configured any
// payment method.
func AnyChannelEnabled() (bool, []PaymentMethod, error) {
	candidates := []string{
		model.OrderPayMethodWechat,
		model.OrderPayMethodAlipay,
		model.OrderPayMethodBank,
	}
	out := make([]PaymentMethod, 0, len(candidates))
	any := false
	for _, name := range candidates {
		if !HasChannel(name) {
			continue
		}
		ch, err := New(name)
		if err != nil {
			continue
		}
		enabled, _ := ch.IsEnabled()
		out = append(out, PaymentMethod{
			Name:    name,
			Label:   PaymentMethodLabel(name),
			Enabled: enabled,
		})
		if enabled {
			any = true
		}
	}
	return any, out, nil
}

// SettingsBool is a small helper that reads a JSON-encoded setting value
// and extracts the boolean "enabled" field. Used by every channel's
// IsEnabled.
func SettingsBool(key string) (bool, error) {
	raw := model.GetSystemSettingString(key)
	if raw == "" {
		return false, nil
	}
	var obj struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return false, fmt.Errorf("解析 %s 失败: %w", key, err)
	}
	return obj.Enabled, nil
}

// SettingsJSON reads a JSON-encoded setting value and unmarshals it
// into out. Returns an error if the value is missing or unparsable.
func SettingsJSON(key string, out interface{}) error {
	raw := model.GetSystemSettingString(key)
	if raw == "" {
		return errors.New("未配置: " + key)
	}
	return json.Unmarshal([]byte(raw), out)
}
