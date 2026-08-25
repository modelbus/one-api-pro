package payment

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat"

	"github.com/modelbus/one-api-pro/model"
)

// wechatChannel implements Channel for WeChat Pay Native (扫码支付).
//
// Credentials are read from the system_settings table at call time:
//   - payment.wechat.enabled -> {"enabled": bool}
//   - payment.wechat.config  -> {"app_id","mch_id","api_key",
//     "notify_url","cert_file","key_file"}
//
// Files referenced by cert_file / key_file must be readable from the
// binary's working directory. Uploaded files are placed under
// data/payment/wechat/ by the settings controller.
type wechatChannel struct{}

func init() {
	RegisterChannel(&wechatChannel{})
}

func (*wechatChannel) Name() string { return model.OrderPayMethodWechat }

func (*wechatChannel) IsEnabled() (bool, error) {
	return SettingsBool(model.SystemSettingKeyWechatEnabled)
}

type wechatConfig struct {
	AppID     string `json:"app_id"`
	MchID     string `json:"mch_id"`
	APIKey    string `json:"api_key"`
	NotifyURL string `json:"notify_url"`
	CertFile  string `json:"cert_file"`
	KeyFile   string `json:"key_file"`
}

func (c *wechatChannel) loadConfig() (*wechatConfig, *wechat.Client, error) {
	var cfg wechatConfig
	if err := SettingsJSON(model.SystemSettingKeyWechatConfig, &cfg); err != nil {
		return nil, nil, fmt.Errorf("未配置微信支付参数: %w", err)
	}
	if cfg.AppID == "" || cfg.MchID == "" || cfg.APIKey == "" {
		return nil, nil, errors.New("微信支付参数不完整 (app_id / mch_id / api_key)")
	}
	// Native prepay only requires app_id / mch_id / api_key. Refund
	// would need the merchant cert+key; refund is out of scope this
	// milestone so the cert paths are not loaded here.
	client := wechat.NewClient(cfg.AppID, cfg.MchID, cfg.APIKey, true)
	return &cfg, client, nil
}

func (c *wechatChannel) PrePay(orderNo string, amount float64, subject string) (*PrePayResult, error) {
	enabled, err := c.IsEnabled()
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, errors.New("微信支付未启用")
	}
	_, client, err := c.loadConfig()
	if err != nil {
		return nil, err
	}

	// amount is yuan; WeChat wants fen (integer).
	totalFen := int64(amount*100 + 0.5)
	if totalFen <= 0 {
		return nil, errors.New("金额必须大于 0")
	}

	bm := gopay.BodyMap{}
	bm.Set("body", subject)
	bm.Set("out_trade_no", orderNo)
	bm.Set("total_fee", totalFen)
	bm.Set("spbill_create_ip", "127.0.0.1")
	bm.Set("notify_url", wechatNotifyURL())
	bm.Set("trade_type", wechat.TradeType_Native)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rsp, err := client.UnifiedOrder(ctx, bm)
	if err != nil {
		return nil, fmt.Errorf("微信下单失败: %w", err)
	}
	if rsp.ReturnCode != "SUCCESS" || rsp.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("微信下单失败: %s/%s", rsp.ReturnCode, rsp.ReturnMsg)
	}
	return &PrePayResult{
		PayURL:   rsp.CodeUrl,
		QRCode:   rsp.CodeUrl,
		ExpireAt: 0, // WeChat Native expiry is merchant-configurable; omitted here
		TradeNo:  rsp.PrepayId,
	}, nil
}

// wechatNotifyXML is the subset of fields we care about.
type wechatNotifyXML struct {
	ReturnCode    string `xml:"return_code"`
	ResultCode    string `xml:"result_code"`
	OutTradeNo    string `xml:"out_trade_no"`
	TransactionID string `xml:"transaction_id"`
	TotalFee      int64  `xml:"total_fee"`
	Sign          string `xml:"sign"`
}

func (c *wechatChannel) VerifyNotify(payload []byte) (*NotifyResult, error) {
	// payload is the raw body of the async POST (XML for WeChat).
	// We use gopay's parser + sign verifier on the parsed struct.
	var n wechatNotifyXML
	if err := xml.Unmarshal(payload, &n); err != nil {
		return nil, fmt.Errorf("解析微信回调 XML 失败: %w", err)
	}
	if !wechatSignValid(&n) {
		return nil, errors.New("微信回调签名校验失败")
	}
	if n.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("微信回调 return_code=%s", n.ReturnCode)
	}
	if n.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("微信回调 result_code=%s", n.ResultCode)
	}
	// total_fee in WeChat is in fen.
	amount := float64(n.TotalFee) / 100.0
	return &NotifyResult{
		OutTradeNo: n.OutTradeNo,
		TradeNo:    n.TransactionID,
		Amount:     amount,
		Paid:       true,
	}, nil
}

// wechatSignValid reconstructs the canonical string and compares MD5
// against the provided sign. Format:
//   k1=v1&k2=v2&...&key=API_KEY
// sorted by ASCII order, excluding `sign`. The API_KEY is appended last
// without a leading "&".
//
// See: https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=4_1
func wechatSignValid(n *wechatNotifyXML) bool {
	key := wechatAPIKey()
	if key == "" {
		return false
	}
	// Re-marshal all known fields except "sign". We use a map for
	// key/value pairs and sort.
	pairs := map[string]string{
		"return_code":    n.ReturnCode,
		"result_code":    n.ResultCode,
		"out_trade_no":   n.OutTradeNo,
		"transaction_id": n.TransactionID,
		"total_fee":      fmt.Sprintf("%d", n.TotalFee),
	}
	keys := make([]string, 0, len(pairs))
	for k := range pairs {
		if pairs[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteString("&")
		}
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(pairs[k])
	}
	b.WriteString("&key=")
	b.WriteString(key)
	sum := md5.Sum([]byte(b.String()))
	got := strings.ToUpper(hex.EncodeToString(sum[:]))
	return got == strings.ToUpper(n.Sign)
}

// wechatAPIKey returns the configured WeChat API key, or "" if absent.
func wechatAPIKey() string {
	var cfg wechatConfig
	_ = SettingsJSON(model.SystemSettingKeyWechatConfig, &cfg)
	return cfg.APIKey
}

// wechatNotifyURL returns the configured notify_url, or "" if absent.
func wechatNotifyURL() string {
	var cfg wechatConfig
	_ = SettingsJSON(model.SystemSettingKeyWechatConfig, &cfg)
	return cfg.NotifyURL
}
