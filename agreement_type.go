package alipay

////////////////////////////////////////////////////////////////////////////////
// https://docs.open.alipay.com/api_2/alipay.user.agreement.query
// 支付宝个人代扣协议查询接口参数
type AliPayAgreementQuery struct {
	PersonalProductCode string `json:"personal_product_code"` // 个人签约产品码，商户和支付宝签约时确定。 必传
	SignScene           string `json:"sign_scene"`            // 协议签约场景，商户和支付宝签约时确定，商户可咨询技术支持。 必传
	AlipayUserId        string `json:"alipay_user_id"`        // 用户的支付宝账号对应 的支付宝唯一用户号，以 2088 开头的 16 位纯数字 组成
	AlipayLogonId       string `json:"alipay_logon_id"`       // 用户的支付宝登录账号，支持邮箱或手机号码格式。本参数与alipay_user_id 不可同时为空，若都填写，则以alipay_user_id 为准
	ExternalAgreementNo string `json:"external_agreement_no"` // 萌推系统内部的 contract_code
}

func (this AliPayAgreementQuery) APIName() string {
	return "alipay.user.agreement.query"
}

func (this AliPayAgreementQuery) Params() map[string]string {
	var m = make(map[string]string)
	return m
}

func (this AliPayAgreementQuery) ExtJSONParamName() string {
	return "biz_content"
}

func (this AliPayAgreementQuery) ExtJSONParamValue() string {
	return marshal(this)
}

type AliPayAgreementQueryResponse struct {
	AlipayUserAgreementQueryResponse struct {
		Code                string `json:"code"`
		Msg                 string `json:"msg"`
		SubCode				string `json:"sub_code"`         	  // 业务错误码
		SubMsg				string `json:"sub_msg"`				  // 业务错误报错信息
		ValidTime           string `json:"valid_time"`            // 协议生效时间，格式为 yyyy-MM-dd HH:mm:ss
		AlipayLogonId       string `json:"alipay_logon_id"`       // 返回脱敏的支付宝账号
		InvalidTime         string `json:"invalid_time"`          // 协议失效时间，格式为 yyyy-MM-dd HH:mm:ss
		PricipalType        string `json:"pricipal_type"`         // 签约主体类型。 CARD:支付宝账号 CUSTOMER:支付宝用户
		DeviceId            string `json:"device_id"`             // 设备id
		PrincipalId         string `json:"principal_id"`          // 签约主体标识。 当principal_type为CARD 时，该字段为支付宝用户号; 当principal_type为 CUSTOMER 时，该字段为支付宝用户标识。一个用户 可能有多个支付宝账号，即多个支付宝用户号，但只有一个是支付宝用户标识。 一个支付宝账号对应一个支付宝唯一用户号(以2088开头的16位纯数字组成)。
		SignScene           string `json:"sign_scene"`            // 签约协议的场景。INDUSTRY|CARRENTAL
		AgreementNo         string `json:"agreement_no"`          // 支付宝签约成功后的协议号
		ThirdPartyType      string `json:"third_party_type"`      // 签约第三方主体类型。对于三方协议，表示当前用户和哪一类的第三方主体进行签约。 取值范围： 1. PARTNER（平台商户）; 2. MERCHANT（集团商户），集团下子商户可共享用户签约内容; 默认为PARTNER
		Status              string `json:"status"`                // 协议当前状态 1. TEMP：暂存，协议未生效过； 2. NORMAL：正常； 3. STOP：暂停
		SignTime            string `json:"sign_time"`             // 协议签约时间，格式为 yyyy-MM-dd HH:mm:ss。
		PersonalProductCode string `json:"personal_product_code"` // 协议产品码，商户和支付宝签约时确定，不同业务场景对应不同的签约产品码
		ExternalAgreementNo string `json:"external_agreement_no"` // 萌推系统内部的 contract_code
		ZmOpenId            string `json:"zm_open_id"`            // 用户的芝麻信用 openId，供商 户查询用户芝麻信用使用。
		ExternalLogonId     string `json:"external_logon_id"`     // 外部登录Id
		CreditAuthMode      string `json:"credit_auth_mode"`      // 授信模式，取值：DEDUCT_HUAZHI-花芝GO。目前只在花芝代扣（即花芝go）协议时才会返回
	} `json:"alipay_user_agreement_query_response"`
	Sign string `json:"sign"`
	Body string `json:"-"` // 返回内容原文，主要是为了记录日志
}

func (this *AliPayAgreementQueryResponse) IsSuccess() bool {
	if this.AlipayUserAgreementQueryResponse.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}
func (this *AliPayAgreementQueryResponse) SetOriginBody(body string) {
	this.Body = body
}

////////////////////////////////////////////////////////////////////////////////
// https://docs.open.alipay.com/api_2/alipay.user.agreement.unsign
// 支付宝个人代扣协议解约接口
type AliPayAgreementUnsign struct {
	NotifyURL           string `json:"-"`                     // 可选
	OperateType         string `json:"operate_type"`          // confirm（解约确认），invalid（解约作废）
	PersonalProductCode string `json:"personal_product_code"` // 个人签约产品码，商户和支付宝签约时确定。 必传
	SignScene           string `json:"sign_scene"`            // 协议签约场景，商户和支付宝签约时确定，商户可咨询技术支持。 必传
	AlipayUserId        string `json:"alipay_user_id"`        // 用户的支付宝账号对应 的支付宝唯一用户号，以 2088 开头的 16 位纯数字 组成
	AlipayLogonId       string `json:"alipay_logon_id"`       // 用户的支付宝登录账号，支持邮箱或手机号码格式。本参数与alipay_user_id 不可同时为空，若都填写，则以alipay_user_id 为准
	ExternalAgreementNo string `json:"external_agreement_no"` // 萌推系统内部的 contract_code
	AgreementNo         string `json:"agreement_no"`          // 支付宝系统中用以唯一标识用户签约记录的编号（用户签约成功后的协议号），如果传了该参数，其他参数会被忽略
}

func (this AliPayAgreementUnsign) APIName() string {
	return "alipay.user.agreement.unsign"
}

func (this AliPayAgreementUnsign) Params() map[string]string {
	var m = make(map[string]string)
	m["notify_url"] = this.NotifyURL
	return m
}

func (this AliPayAgreementUnsign) ExtJSONParamName() string {
	return "biz_content"
}

func (this AliPayAgreementUnsign) ExtJSONParamValue() string {
	return marshal(this)
}

type AliPayAgreementUnsignResponse struct {
	AlipayUserAgreementUnsignResponse struct {
		Code    string `json:"code"`
		Msg     string `json:"msg"`
		SubCode string `json:"sub_code"`
		SubMsg  string `json:"sub_msg"`
	} `json:"alipay_user_agreement_unsign_response"`
	Sign string `json:"sign"`
	Body string `json:"-"` // 返回内容原文，主要是为了记录日志
}

func (this *AliPayAgreementUnsignResponse) IsSuccess() bool {
	if this.AlipayUserAgreementUnsignResponse.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}
func (this *AliPayAgreementUnsignResponse) SetOriginBody(body string) {
	this.Body = body
}

////////////////////////////////////////////////////////////////////////////////
// https://docs.open.alipay.com/api_2/alipay.user.agreement.executionplan.modify
// 支付宝个人代扣协议修改接口
type AliPayAgreementModify struct {
	AgreementNo string `json:"agreement_no"` // 支付宝系统中用以唯一标识用户签约记录的编号（用户签约成功后的协议号），如果传了该参数，其他参数会被忽略
	DeductTime  string `json:"deduct_time"`  // 商户下一次扣款时间
	Memo        string `json:"memo"`         // 具体修改原因
}

func (this AliPayAgreementModify) APIName() string {
	return "alipay.user.agreement.executionplan.modify"
}

func (this AliPayAgreementModify) Params() map[string]string {
	var m = make(map[string]string)
	return m
}

func (this AliPayAgreementModify) ExtJSONParamName() string {
	return "biz_content"
}

func (this AliPayAgreementModify) ExtJSONParamValue() string {
	return marshal(this)
}

type AliPayAgreementModifyResponse struct {
	AlipayUserAgreementExecutionplanModifyResponse struct {
		Code    	string `json:"code"`
		Msg    	 	string `json:"msg"`
		SubCode 	string `json:"sub_code"`
		SubMsg  	string `json:"sub_msg"`
		AgreementNo string `json:"agreement_no"`
		DeductTime 	string `json:"deduct_time"`
	} `json:"alipay_user_agreement_executionplan_modify_response"`
	Sign string `json:"sign"`
	Body string `json:"-"` // 返回内容原文，主要是为了记录日志
}

func (this *AliPayAgreementModifyResponse) IsSuccess() bool {
	if this.AlipayUserAgreementExecutionplanModifyResponse.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}
func (this *AliPayAgreementModifyResponse) SetOriginBody(body string) {
	this.Body = body
}