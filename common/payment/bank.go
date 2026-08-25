package payment

import (
	"errors"

	"github.com/modelbus/one-api-pro/model"
)

// bankChannel implements Channel for "bank" / "offline" payments.
// These are placeholder channels — no upstream SDK is invoked. The
// order stays in status=0 (pending) until an admin marks it paid
// via the admin order endpoints (PUT /api/order/:id).
type bankChannel struct{}

func init() {
	RegisterChannel(&bankChannel{})
	RegisterChannel(&offlineChannel{})
	RegisterChannel(&freeChannel{})
}

func (*bankChannel) Name() string { return model.OrderPayMethodBank }

func (*bankChannel) IsEnabled() (bool, error) {
	return SettingsBool(model.SystemSettingKeyBankEnabled)
}

func (*bankChannel) PrePay(_ string, _ float64, _ string) (*PrePayResult, error) {
	return nil, errors.New("bank 支付未实现：等待管理员在后台标记收款")
}

func (*bankChannel) VerifyNotify(_ []byte) (*NotifyResult, error) {
	return nil, errors.New("bank 支付没有异步回调")
}

// offlineChannel is identical to bank for now.
type offlineChannel struct{}

func (*offlineChannel) Name() string { return model.OrderPayMethodOffline }

func (*offlineChannel) IsEnabled() (bool, error) {
	// Always considered enabled; if disabled the API path is not used.
	return true, nil
}

func (*offlineChannel) PrePay(_ string, _ float64, _ string) (*PrePayResult, error) {
	return nil, errors.New("offline 支付未实现：等待管理员在后台标记收款")
}

func (*offlineChannel) VerifyNotify(_ []byte) (*NotifyResult, error) {
	return nil, errors.New("offline 支付没有异步回调")
}

// freeChannel is the "admin free grant" path — PrePay is unused because
// the admin open-subscription flow creates the order already paid.
type freeChannel struct{}

func (*freeChannel) Name() string { return model.OrderPayMethodFree }

func (*freeChannel) IsEnabled() (bool, error) { return true, nil }

func (*freeChannel) PrePay(_ string, _ float64, _ string) (*PrePayResult, error) {
	return &PrePayResult{PayURL: "", QRCode: "", ExpireAt: 0, TradeNo: ""}, nil
}

func (*freeChannel) VerifyNotify(_ []byte) (*NotifyResult, error) {
	return &NotifyResult{Paid: true}, nil
}
