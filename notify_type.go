package alipay

import (
	"errors"
	"net/http"
)

/*
	支付：notify_type=trade_status_sync&trade_status=TRADE_SUCCESS
	签约：notify_type=dut_user_sign&status=NORMAL
	解约：notify_type=dut_user_unsign&status=UNSIGN
	支付和签约可以通过这些条件进行判断，notify_type不同，并且有另一个字段来标记 支付的状态 和 签约的状态
*/

// https://doc.open.alipay.com/docs/doc.htm?spm=a219a.7629140.0.0.8AmJwg&treeId=203&articleId=105286&docType=1
type TradeNotification struct {
	AuthAppId         string `json:"auth_app_id"`         // App Id
	NotifyTime        string `json:"notify_time"`         // 通知时间
	NotifyType        string `json:"notify_type"`         // 通知类型
	NotifyId          string `json:"notify_id"`           // 通知校验ID
	AppId             string `json:"app_id"`              // 开发者的app_id
	Charset           string `json:"charset"`             // 编码格式
	Version           string `json:"version"`             // 接口版本
	SignType          string `json:"sign_type"`           // 签名类型
	Sign              string `json:"sign"`                // 签名
	TradeNo           string `json:"trade_no"`            // 支付宝交易号
	OutTradeNo        string `json:"out_trade_no"`        // 商户订单号
	OutBizNo          string `json:"out_biz_no"`          // 商户业务号
	BuyerId           string `json:"buyer_id"`            // 买家支付宝用户号
	BuyerLogonId      string `json:"buyer_logon_id"`      // 买家支付宝账号
	SellerId          string `json:"seller_id"`           // 卖家支付宝用户号
	SellerEmail       string `json:"seller_email"`        // 卖家支付宝账号
	TradeStatus       string `json:"trade_status"`        // 交易状态
	TotalAmount       string `json:"total_amount"`        // 订单金额
	ReceiptAmount     string `json:"receipt_amount"`      // 实收金额
	InvoiceAmount     string `json:"invoice_amount"`      // 开票金额
	BuyerPayAmount    string `json:"buyer_pay_amount"`    // 付款金额
	PointAmount       string `json:"point_amount"`        // 集分宝金额
	RefundFee         string `json:"refund_fee"`          // 总退款金额
	Subject           string `json:"subject"`             // 总退款金额
	Body              string `json:"body"`                // 商品描述
	GmtCreate         string `json:"gmt_create"`          // 交易创建时间
	GmtPayment        string `json:"gmt_payment"`         // 交易付款时间
	GmtRefund         string `json:"gmt_refund"`          // 交易退款时间
	GmtClose          string `json:"gmt_close"`           // 交易结束时间
	FundBillList      string `json:"fund_bill_list"`      // 支付金额信息
	PassbackParams    string `json:"passback_params"`     // 回传参数
	VoucherDetailList string `json:"voucher_detail_list"` // 优惠券信息
}

func (this *AliPay) GetTradeNotification(req *http.Request) (*TradeNotification, error) {
	return GetTradeNotification(req, this.AliPayPublicKey)
}

func GetTradeNotification(req *http.Request, aliPayPublicKey []byte) (noti *TradeNotification, err error) {
	if req == nil {
		return nil, errors.New("request 参数不能为空")
	}
	req.ParseForm()

	noti = &TradeNotification{}
	noti.AppId = req.FormValue("app_id")
	noti.AuthAppId = req.FormValue("auth_app_id")
	noti.NotifyId = req.FormValue("notify_id")
	noti.NotifyType = req.FormValue("notify_type")
	noti.NotifyTime = req.FormValue("notify_time")
	noti.TradeNo = req.FormValue("trade_no")
	noti.TradeStatus = req.FormValue("trade_status")
	noti.TotalAmount = req.FormValue("total_amount")
	noti.ReceiptAmount = req.FormValue("receipt_amount")
	noti.InvoiceAmount = req.FormValue("invoice_amount")
	noti.BuyerPayAmount = req.FormValue("buyer_pay_amount")
	noti.SellerId = req.FormValue("seller_id")
	noti.SellerEmail = req.FormValue("seller_email")
	noti.BuyerId = req.FormValue("buyer_id")
	noti.BuyerLogonId = req.FormValue("buyer_logon_id")
	noti.FundBillList = req.FormValue("fund_bill_list")
	noti.Charset = req.FormValue("charset")
	noti.PointAmount = req.FormValue("point_amount")
	noti.OutTradeNo = req.FormValue("out_trade_no")
	noti.OutBizNo = req.FormValue("out_biz_no")
	noti.GmtCreate = req.FormValue("gmt_create")
	noti.GmtPayment = req.FormValue("gmt_payment")
	noti.GmtRefund = req.FormValue("gmt_refund")
	noti.GmtClose = req.FormValue("gmt_close")
	noti.Subject = req.FormValue("subject")
	noti.Body = req.FormValue("body")
	noti.RefundFee = req.FormValue("refund_fee")
	noti.Version = req.FormValue("version")
	noti.SignType = req.FormValue("sign_type")
	noti.Sign = req.FormValue("sign")
	noti.PassbackParams = req.FormValue("passback_params")
	noti.VoucherDetailList = req.FormValue("voucher_detail_list")

	if len(noti.NotifyId) == 0 {
		return nil, errors.New("不是有效的 Notify")
	}

	ok, err := verifySign(req, aliPayPublicKey)
	if ok {
		return noti, nil
	}
	return nil, err
}

// https://docs.open.alipay.com/api_2/alipay.user.agreement.page.sign
type AgreementSignNotification struct {
	NotifyType          string `json:"notify_type"`           // 通知类型
	NotifyId            string `json:"notify_id"`             // 通知校验ID
	NotifyTime          string `json:"notify_time"`           // 通知时间
	SignType            string `json:"sign_type"`             // 签名类型
	Sign                string `json:"sign"`                  // 签名
	AppId               string `json:"app_id"`                // 开发者的app_id
	AuthAppId           string `json:"auth_app_id"`           // App Id
	SignTime            string `json:"sign_time"`             // 签约时间
	ExternalLogonId     string `json:"external_logon_id"`     // 用户在商户网站的登录账号，如果商户接口中未传，则不会返回
	ExternalAgreementNo string `json:"external_agreement_no"` // 萌推系统内部的 contract_code
	PersonalProductCode string `json:"personal_product_code"` // 个人签约产品码，商户和支付宝签约时确定。 必传
	SignScene           string `json:"sign_scene"`            // 协议签约场景
	AlipayUserId        string `json:"alipay_user_id"`        // 支付宝用户唯一id
	AlipayLogonId       string `json:"alipay_logon_id"`       // 用户的支付宝登录账号，支持邮箱或手机号码格式。本参数与alipay_user_id 不可同时为空，若都填写，则以alipay_user_id 为准
	Status              string `json:"status"`                // 协议的当前状态。 1. TEMP：暂存，协议未生效过； 2. NORMAL：正常； 3. STOP：暂停。 4. UNSIGN：解约。（只有签约成功才会返回）
	ValidTime           string `json:"valid_time"`            // 用户代扣协议的实际生效时间，格式为yyyy-MM-dd HH:mm:ss。（只有签约成功才会返回）
	InvalidTime         string `json:"invalid_time"`          // 用户代扣协议的失效时间，格式为yyyy-MM-dd HH:mm:ss。（只有签约成功才会返回）
	AgreementNo         string `json:"agreement_no"`          // 支付宝系统中用以唯一标识用户签约记录的编号。（只有签约成功时才会返回）
	ZmOpenId            string `json:"zm_open_id"`            // 用户的芝麻信用openId，供商户查询用户芝麻信用使用。（只有签约成功时才返回）
	ForexEligible       string `json:"forex_eligible"`        // 是否海外购汇身份。值：T/F（只有在签约成功时才会返回）
	UnsignTime          string `json:"unsign_time"`           // 解约时间
}

func (this *AliPay) GetAgreementSignNotification(req *http.Request) (*AgreementSignNotification, error) {
	return GetAgreementSignNotification(req, this.AliPayPublicKey)
}

func GetAgreementSignNotification(req *http.Request, aliPayPublicKey []byte) (noti *AgreementSignNotification, err error) {
	if req == nil {
		return nil, errors.New("request 参数不能为空")
	}
	req.ParseForm()

	noti = &AgreementSignNotification{}
	noti.NotifyType = req.FormValue("notify_type")
	noti.NotifyId = req.FormValue("notify_id")
	noti.NotifyTime = req.FormValue("notify_time")
	noti.SignType = req.FormValue("sign_type")
	noti.Sign = req.FormValue("sign")
	noti.AppId = req.FormValue("app_id")
	noti.AuthAppId = req.FormValue("auth_app_id")
	noti.SignTime = req.FormValue("sign_time")
	noti.ExternalLogonId = req.FormValue("external_logon_id")
	noti.ExternalAgreementNo = req.FormValue("external_agreement_no")
	noti.PersonalProductCode = req.FormValue("personal_product_code")
	noti.SignScene = req.FormValue("sign_scene")
	noti.UnsignTime = req.FormValue("unsign_time")
	noti.AlipayUserId = req.FormValue("alipay_user_id")
	noti.AlipayLogonId = req.FormValue("alipay_logon_id")
	noti.Status = req.FormValue("status")
	noti.ValidTime = req.FormValue("valid_time")
	noti.InvalidTime = req.FormValue("invalid_time")
	noti.AgreementNo = req.FormValue("agreement_no")
	noti.ZmOpenId = req.FormValue("zm_open_id")
	noti.ForexEligible = req.FormValue("forex_eligible")
	if len(noti.NotifyId) == 0 {
		return nil, errors.New("不是有效的 SignNotify")
	}
	ok, err := verifySign(req, aliPayPublicKey)
	if ok {
		return noti, nil
	}
	return nil, err
}
