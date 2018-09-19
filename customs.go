package alipay

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// https://docs.open.alipay.com/155/104778/
type DeclareRequest struct {
	// 必选
	OutRequestNo        string `xml:"out_request_no"`
	TradeNo             string `xml:"trade_no"`
	MerchantCustomsCode string `xml:"merchant_customs_code"`
	MerchantCustomsName string `xml:"merchant_customs_name"`
	Amount              string `xml:"amount"`
	CustomsPlace        string `xml:"customs_place"`

	// 可选
	IsSplit     string `xml:"is_split"`
	SubOutBizNo string `xml:"sub_out_biz_no"`
	BuyerName   string `xml:"buyer_name"`
	BuyerIdNo   string `xml:"buyer_id_no"`
}

func (this DeclareRequest) APIName() string {
	return "alipay.acquire.customs"
}

func (this DeclareRequest) Params() map[string]string {
	var m = make(map[string]string)
	//m["partner"] = this.Partner
	m["_input_charset"] = "UTF-8"
	m["out_request_no"] = this.OutRequestNo
	m["trade_no"] = this.TradeNo
	m["merchant_customs_code"] = this.MerchantCustomsCode
	m["merchant_customs_name"] = this.MerchantCustomsName
	m["amount"] = this.Amount
	m["customs_place"] = this.CustomsPlace
	m["service"] = this.APIName()

	// 可选
	if this.IsSplit != "" {
		m["is_split"] = this.IsSplit
		m["sub_out_biz_no"] = this.SubOutBizNo
		m["buyer_name"] = this.BuyerName
		m["buyer_id_no"] = this.BuyerIdNo
	}
	return m
}

type DeclareResponse struct {
	Body string `xml:"-"` // 方便记日志用

	IsSuccess string `xml:"is_success"`
	SignType  string `xml:"sign_type"`
	Sign      string `xml:"sign"`
	Error     string `xml:"error"`

	ResultCode      string `xml:"result_code"`
	TradeNo         string `xml:"trade_no"`
	AlipayDeclareNo string `xml:"alipay_declare_no"`
	DetailErrorCode string `xml:"detail_error_code"`
	DetailErrorDes  string `xml:"detail_error_des"`

	IdentityCheck string `xml:"identity_check"`
	// 订购人身份信息和支付人身份信息的验证结果
	// T表示商户传入的订购人姓名和身份证号和支付人的姓名和身份证号一致。
	// F代表商户传入的订购人姓名和身份证号和支付人的姓名和身份证号不一致。
	// 对于同一笔报关单支付宝只会校验一次，如果多次重推不会返回此参数。
}

func (this *DeclareResponse) SetOriginBody(body string) {
	this.Body = body
}

func (this *AliPay) CustomsDeclare(param DeclareRequest) (result *DeclareResponse, err error) {
	// 老接口，所以单独写一套
	var p = url.Values{}
	p.Add("partner", this.partnerId)
	p.Add("charset", K_CHARSET)
	p.Add("sign_type", this.SignType)

	var ps = param.Params()
	if ps != nil {
		for key, value := range ps {
			p.Add(key, value)
		}
	}

	var keys = make([]string, 0, 0)
	for key, _ := range p {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var sign string
	if this.SignType == K_SIGN_TYPE_RSA {
		sign, err = signRSA(keys, p, this.privateKey)
	} else {
		sign, err = signRSA2(keys, p, this.privateKey)
	}
	if err != nil {
		return nil, err
	}
	p.Add("sign", sign)

	buf := strings.NewReader(p.Encode())

	req, err := http.NewRequest(http.MethodGet, this.apiDomain, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := this.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result = &DeclareResponse{}
	err = xml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	result.SetOriginBody(string(data))

	if len(this.AliPayPublicKey) > 0 {
		if ok, err := verifyResponseData(data, result.SignType, result.Sign, this.AliPayPublicKey); ok == false {
			return nil, err
		}
	}

	return
}
