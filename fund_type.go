package alipay

////////////////////////////////////////////////////////////////////////////////
// https://doc.open.alipay.com/docs/api.htm?apiId=1321&docType=4
// 单笔转账到支付宝账户接口请求结构
type AliPayFundTransToAccountTransfer struct {
	AppAuthToken  string `json:"-"`               // 可选
	OutBizNo      string `json:"out_biz_no"`      // 必选 商户转账唯一订单号
	PayeeType     string `json:"payee_type"`      // 必选 收款方账户类型,"ALIPAY_LOGONID":支付宝帐号
	PayeeAccount  string `json:"payee_account"`   // 必选 收款方账户。与payee_type配合使用
	Amount        string `json:"amount"`          // 必选 转账金额,元
	PayerShowName string `json:"payer_show_name"` // 可选 付款方显示姓名
	PayeeRealName string `json:"payee_real_name"` // 可选 收款方真实姓名,如果本参数不为空，则会校验该账户在支付宝登记的实名是否与收款方真实姓名一致。
	Remark        string `json:"remark"`          // 可选 转账备注,金额大于50000时必填
}

func (this AliPayFundTransToAccountTransfer) APIName() string {
	return "alipay.fund.trans.toaccount.transfer"
}

func (this AliPayFundTransToAccountTransfer) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

func (this AliPayFundTransToAccountTransfer) ExtJSONParamName() string {
	return "biz_content"
}

func (this AliPayFundTransToAccountTransfer) ExtJSONParamValue() string {
	return marshal(this)
}

// 单笔转账到支付宝账户接口响应参数
type AliPayFundTransToAccountTransferResponse struct {
	Body struct {
		Code     string `json:"code"`
		Msg      string `json:"msg"`
		SubCode  string `json:"sub_code"`
		SubMsg   string `json:"sub_msg"`
		OutBizNo string `json:"out_biz_no"` // 商户转账唯一订单号：发起转账来源方定义的转账单据号。请求时对应的参数，原样返回
		OrderId  string `json:"order_id"`   // 支付宝转账单据号，成功一定返回，失败可能不返回也可能返回
		PayDate  string `json:"pay_date"`   // 支付时间：格式为yyyy-MM-dd HH:mm:ss，仅转账成功返回
	} `json:"alipay_fund_trans_toaccount_transfer_response"`
	Sign string `json:"sign"`
}

func (this *AliPayFundTransToAccountTransferResponse) IsSuccess() bool {
	if this.Body.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// https://doc.open.alipay.com/docs/api.htm?spm=a219a.7395905.0.0.SIkNrH&docType=4&apiId=1322
// 查询转账订单接口请求参数
type AliPayFundTransOrderQuery struct {
	AppAuthToken string `json:"-"`                    // 可选
	OutBizNo     string `json:"out_biz_no,omitempty"` // 与 OrderId 二选一
	OrderId      string `json:"order_id,omitempty"`   // 与 OutBizNo 二选一
}

func (this AliPayFundTransOrderQuery) APIName() string {
	return "alipay.fund.trans.order.query"
}

func (this AliPayFundTransOrderQuery) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

func (this AliPayFundTransOrderQuery) ExtJSONParamName() string {
	return "biz_content"
}

func (this AliPayFundTransOrderQuery) ExtJSONParamValue() string {
	return marshal(this)
}

// 查询转账订单接口响应参数
type AliPayFundTransOrderQueryResponse struct {
	Body struct {
		Code           string `json:"code"`
		Msg            string `json:"msg"`
		SubCode        string `json:"sub_code"`
		SubMsg         string `json:"sub_msg"`
		OutBizNo       string `json:"out_biz_no"`       // 发起转账来源方定义的转账单据号。 该参数的赋值均以查询结果中 的 out_biz_no 为准。 如果查询失败，不返回该参数
		OrderId        string `json:"order_id"`         // 支付宝转账单据号，查询失败不返回。
		Status         string `json:"status"`           // 转账单据状态
		PayDate        string `json:"pay_date"`         // 支付时间
		ArrivalTimeEnd string `json:"arrival_time_end"` // 预计到账时间，转账到银行卡专用
		OrderFree      string `json:"order_fee"`        // 预计收费金额（元），转账到银行卡专用
		FailReason     string `json:"fail_reason"`      // 查询到的订单状态为FAIL失败或REFUND退票时，返回具体的原因。
		ErrorCode      string `json:"error_code"`       // 查询失败时，本参数为错误代 码。 查询成功不返回。 对于退票订单，不返回该参数。
	} `json:"alipay_fund_trans_order_query_response"`
	Sign string `json:"sign"`
}

func (this *AliPayFundTransOrderQueryResponse) IsSuccess() bool {
	if this.Body.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// https://opendocs.alipay.com/apis/api_28/alipay.fund.trans.uni.transfer/
// 单笔转账接口（新）请求结构
type TransferParticipant struct {
	Identity     string `json:"identity"`      // 必填 参与方的唯一标识
	IdentityType string `json:"identity_type"` // 必填 参与方的标识类型，目前支持如下类型：1、ALIPAY_USER_ID 支付宝的会员ID 2、ALIPAY_LOGON_ID：支付宝登录号，支持邮箱和手机号格式
	Name         string `json:"name"`          // 可选 参与方真实姓名，如果非空，将校验收款支付宝账号姓名一致性。当identity_type=ALIPAY_LOGON_ID时，本字段必填。
}
type FundTransUniTransfer struct {
	AppAuthToken    string               `json:"-"`                 // 可选
	OutBizNo        string               `json:"out_biz_no"`        // 必选 商户转账唯一订单号
	TransAmount     string               `json:"trans_amount"`      // 必选 订单总金额，单位为元
	ProductCode     string               `json:"product_code"`      // 必选 业务产品码，单笔无密转账到支付宝账户固定为: TRANS_ACCOUNT_NO_PWD； 单笔无密转账到银行卡固定为: RANS_BANKCARD_NO_PWD; 发现金红包固定为: TD_RED_PACKET；
	BizScene        string               `json:"biz_scene"`         // 可选 描述特定的业务场景,DIRECT_TRANSFER：单笔无密转账到支付宝/银行卡, B2C现金红包;PERSONAL_COLLECTION：C2C现金红包-领红包
	OrderTitle      string               `json:"order_title"`       // 可选 转账业务的标题，用于在支付宝用户的账单里显示
	OriginalOrderId string               `json:"original_order_id"` // 可选 原支付宝业务单号。C2C现金红包-红包领取时，传红包支付时返回的支付宝单号；B2C现金红包、单笔无密转账到支付宝/银行卡不需要该参数。
	PayeeInfo       *TransferParticipant `json:"payee_info"`        // 必选 收款方信息
	Remark          string               `json:"remark"`            // 可选 业务备注
	BusinessParams  string               `json:"business_params"`   // 可选 转账业务请求的扩展参数，支持传入的扩展参数如下： 1、sub_biz_scene 子业务场景，红包业务必传，取值REDPACKET，C2C现金红包、B2C现金红包均需传入；2、withdraw_timeliness为转账到银行卡的预期到账时间，可选（不传入则默认为T1），取值T0表示预期T+0到账，取值T1表示预期T+1到账，因到账时效受银行机构处理影响，支付宝无法保证一定是T0或者T1到账；

}

func (this FundTransUniTransfer) APIName() string {
	return "alipay.fund.trans.uni.transfer"
}

func (this FundTransUniTransfer) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}

func (this FundTransUniTransfer) ExtJSONParamName() string {
	return "biz_content"
}

func (this FundTransUniTransfer) ExtJSONParamValue() string {
	return marshal(this)
}

// 单笔转账接口(新)响应参数
type FundTransUniTransferResponse struct {
	Body struct {
		Code           string `json:"code"`
		Msg            string `json:"msg"`
		SubCode        string `json:"sub_code"`
		SubMsg         string `json:"sub_msg"`
		OutBizNo       string `json:"out_biz_no"`        // 商户转账唯一订单号：发起转账来源方定义的转账单据号。请求时对应的参数，原样返回
		OrderId        string `json:"order_id"`          // 支付宝转账单据号，成功一定返回，失败可能不返回也可能返回
		PayFundOrderId string `json:"pay_fund_order_id"` // 支付宝支付资金流水号
		Status         string `json:"status"`            // 转账单据状态。
		TransDate      string `json:"trans_date"`        // 订单支付时间
	} `json:"alipay_fund_trans_uni_transfer_response"`
	Sign string `json:"sign"`
}

func (this *FundTransUniTransferResponse) IsSuccess() bool {
	if this.Body.Code == K_SUCCESS_CODE {
		return true
	}
	return false
}
