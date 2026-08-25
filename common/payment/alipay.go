package payment

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"

	"github.com/modelbus/one-api-pro/model"
)

// alipayChannel implements Channel for Alipay Face-to-Face (当面付) QR.
//
// Credentials are read from the system_settings table at call time:
//   - payment.alipay.enabled -> {"enabled": bool}
//   - payment.alipay.config  -> {"app_id","private_key","public_key",
//     "private_key_file","public_key_file","notify_url","gateway"}
//
// Files referenced by *_key_file are PEM-encoded RSA keys. Private
// key signs outgoing prepay; Alipay public key verifies incoming
// notifications. Both are uploaded by the settings controller to
// data/payment/alipay/.
type alipayChannel struct{}

func init() {
	RegisterChannel(&alipayChannel{})
}

func (*alipayChannel) Name() string { return model.OrderPayMethodAlipay }

func (*alipayChannel) IsEnabled() (bool, error) {
	return SettingsBool(model.SystemSettingKeyAlipayEnabled)
}

type alipayConfig struct {
	AppID          string `json:"app_id"`
	PrivateKey     string `json:"private_key"`     // inline key, optional
	PublicKey      string `json:"public_key"`      // inline key, optional
	PrivateKeyFile string `json:"private_key_file"` // path to private key PEM file (preferred)
	PublicKeyFile  string `json:"public_key_file"`  // path to public key PEM file (preferred)
	NotifyURL      string `json:"notify_url"`
	Gateway        string `json:"gateway"` // https://openapi.alipay.com/gateway.do (prod) or sandbox
}

// alipayNotifyBean is the subset of fields used for sign verification.
// The fields are tagged with "json" because alipay.VerifySign uses
// json.Marshal to produce the canonical sign string.
type alipayNotifyBean struct {
	NotifyTime    string `json:"notify_time"`
	NotifyType    string `json:"notify_type"`
	NotifyID      string `json:"notify_id"`
	AppID         string `json:"app_id"`
	OutTradeNo    string `json:"out_trade_no"`
	OutBizNo      string `json:"out_biz_no"`
	TradeNo       string `json:"trade_no"`
	TradeStatus   string `json:"trade_status"`
	TotalAmount   string `json:"total_amount"`
	ReceiptAmount string `json:"receipt_amount"`
	BuyerID       string `json:"buyer_id"`
	SellerID      string `json:"seller_id"`
	Subject       string `json:"subject"`
	Body          string `json:"body"`
	GmtCreate     string `json:"gmt_create"`
	GmtPayment    string `json:"gmt_payment"`
	GmtClose      string `json:"gmt_close"`
	FundBillList  string `json:"fund_bill_list"`
	PassbackParams string `json:"passback_params"`
	Version       string `json:"version"`
	Charset       string `json:"charset"`
}

func (c *alipayChannel) loadConfig() (*alipayConfig, error) {
	var cfg alipayConfig
	if err := SettingsJSON(model.SystemSettingKeyAlipayConfig, &cfg); err != nil {
		return nil, fmt.Errorf("未配置支付宝参数: %w", err)
	}
	if cfg.AppID == "" {
		return nil, errors.New("支付宝参数不完整 (app_id)")
	}
	// File-based keys take precedence; otherwise fall back to inline.
	if cfg.PrivateKey == "" && cfg.PrivateKeyFile != "" {
		if data, err := os.ReadFile(cfg.PrivateKeyFile); err == nil {
			cfg.PrivateKey = string(data)
		}
	}
	if cfg.PublicKey == "" && cfg.PublicKeyFile != "" {
		if data, err := os.ReadFile(cfg.PublicKeyFile); err == nil {
			cfg.PublicKey = string(data)
		}
	}
	if cfg.PrivateKey == "" {
		return nil, errors.New("支付宝参数不完整 (private_key)")
	}
	if cfg.Gateway == "" {
		cfg.Gateway = "https://openapi.alipay.com/gateway.do"
	}
	return &cfg, nil
}

func (c *alipayChannel) client() (*alipay.Client, *alipayConfig, error) {
	cfg, err := c.loadConfig()
	if err != nil {
		return nil, nil, err
	}
	client, err := alipay.NewClient(cfg.AppID, cfg.PrivateKey, true)
	if err != nil {
		return nil, nil, fmt.Errorf("创建支付宝客户端失败: %w", err)
	}
	if cfg.PublicKey != "" {
		client.AutoVerifySign([]byte(cfg.PublicKey))
	}
	return client, cfg, nil
}

func (c *alipayChannel) PrePay(orderNo string, amount float64, subject string) (*PrePayResult, error) {
	enabled, err := c.IsEnabled()
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, errors.New("支付宝未启用")
	}
	client, cfg, err := c.client()
	if err != nil {
		return nil, err
	}
	_ = cfg // reserved for future use

	amountStr := fmt.Sprintf("%.2f", amount)

	bm := gopay.BodyMap{}
	bm.Set("out_trade_no", orderNo)
	bm.Set("total_amount", amountStr)
	bm.Set("subject", subject)
	bm.Set("notify_url", alipayNotifyURL())

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rsp, err := client.TradePrecreate(ctx, bm)
	if err != nil {
		return nil, fmt.Errorf("支付宝预下单失败: %w", err)
	}
	if rsp == nil || rsp.Response == nil || rsp.Response.Code != "10000" {
		code := ""
		msg := ""
		if rsp != nil && rsp.Response != nil {
			code = rsp.Response.Code
			msg = rsp.Response.Msg
		}
		return nil, fmt.Errorf("支付宝预下单失败: %s/%s", code, msg)
	}
	// For precreate, Alipay returns a qr_code string the client should
	// render as a QR for the user to scan.
	qr := rsp.Response.QrCode
	return &PrePayResult{
		PayURL:   qr,
		QRCode:   qr,
		ExpireAt: 0, // Alipay default 15min; merchants can shorten via timeout_express
		TradeNo:  "",
	}, nil
}

func (c *alipayChannel) VerifyNotify(payload []byte) (*NotifyResult, error) {
	cfg, err := c.loadConfig()
	if err != nil {
		return nil, err
	}
	if cfg.PublicKey == "" {
		return nil, errors.New("支付宝公钥未配置，无法验签")
	}
	values, err := parseAlipayForm(string(payload))
	if err != nil {
		return nil, err
	}
	bean := alipayFormToBean(values)
	sign := values["sign"]
	if sign == "" {
		return nil, errors.New("支付宝回调缺少 sign 字段")
	}
	ok, err := alipay.VerifySign(cfg.PublicKey, bean)
	if err != nil {
		return nil, fmt.Errorf("支付宝验签错误: %w", err)
	}
	if !ok {
		return nil, errors.New("支付宝回调签名校验失败")
	}
	tradeStatus := values["trade_status"]
	paid := tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED"
	if !paid {
		return nil, fmt.Errorf("支付宝回调 trade_status=%s", tradeStatus)
	}
	amount, err := parseAlipayAmount(values["total_amount"])
	if err != nil {
		return nil, err
	}
	return &NotifyResult{
		OutTradeNo: values["out_trade_no"],
		TradeNo:    values["trade_no"],
		Amount:     amount,
		Paid:       true,
	}, nil
}

// parseAlipayForm parses an x-www-form-urlencoded body into a map.
func parseAlipayForm(body string) (map[string]string, error) {
	raw := map[string]string{}
	if body == "" {
		return raw, errors.New("支付宝回调 body 为空")
	}
	for _, pair := range strings.Split(body, "&") {
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		k, err := url.QueryUnescape(kv[0])
		if err != nil {
			return nil, err
		}
		v, err := url.QueryUnescape(kv[1])
		if err != nil {
			return nil, err
		}
		raw[k] = v
	}
	return raw, nil
}

// alipayFormToBean maps the form-decoded key/value pairs to a
// alipayNotifyBean suitable for alipay.VerifySign (which uses
// json.Marshal to produce the canonical sign string).
func alipayFormToBean(v map[string]string) alipayNotifyBean {
	return alipayNotifyBean{
		NotifyTime:     v["notify_time"],
		NotifyType:     v["notify_type"],
		NotifyID:       v["notify_id"],
		AppID:          v["app_id"],
		OutTradeNo:     v["out_trade_no"],
		OutBizNo:       v["out_biz_no"],
		TradeNo:        v["trade_no"],
		TradeStatus:    v["trade_status"],
		TotalAmount:    v["total_amount"],
		ReceiptAmount:  v["receipt_amount"],
		BuyerID:        v["buyer_id"],
		SellerID:       v["seller_id"],
		Subject:        v["subject"],
		Body:           v["body"],
		GmtCreate:      v["gmt_create"],
		GmtPayment:     v["gmt_payment"],
		GmtClose:       v["gmt_close"],
		FundBillList:   v["fund_bill_list"],
		PassbackParams: v["passback_params"],
		Version:        v["version"],
		Charset:        v["charset"],
	}
}

func parseAlipayAmount(s string) (float64, error) {
	if s == "" {
		return 0, errors.New("支付宝回调金额为空")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("解析支付宝金额 %q 失败: %w", s, err)
	}
	return f, nil
}

// alipayNotifyURL returns the configured notify_url, or "" if absent.
func alipayNotifyURL() string {
	var cfg alipayConfig
	_ = SettingsJSON(model.SystemSettingKeyAlipayConfig, &cfg)
	return cfg.NotifyURL
}
