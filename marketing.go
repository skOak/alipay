package alipay

import (
    "encoding/json"
)

func (this *AliPay) MarketingConsult( param MarketingConsultParam) (*MarketingConsultResponse, error) {
    results := &MarketingConsultResponse{}
    err := this.doRequest("POST", param, results)
    return results, err
}

// marketing consult
//type MarketingConsultReq struct {
//    BizScene string
//}
//
type MarketingConsultResponse struct {
   Text        string //是 2048 营销文案文本内容 使用支付宝支付笔笔享优惠
   Image       string //否 2048 营销图片地址 http://www.xxx.com
   PrePayToken string //否 64 营销联动唯一标识，商户在后续调收单接口时，需将此字段 透传带入 xxxxxxx
   ExtInfo     string //否 2048 扩展信息，json格式 {"icon":"http:xxx"}
}

type MarketingConsultParam struct {
    BizScene        string `json:"biz_scene"`   //64 业务场景，用于区分商户具体的咨场景；OPENING_PAGE：开屏页营销咨询；DETAIL_PAGE：详情页营销咨询
    Mobile          string `json:"mobile"`
    EncryptedMobile string `json:"encrypted_mobile"`

    //AppAuthToken string `json:"app_auth_token"`
    OutTradeNo string `json:"out_trade_no"`
    TotalAmount string `json:"total_amount"`
    UndiscountableAmount string `json:"undiscountable_amount"`
}

func (this MarketingConsultParam) APIName() string {
    return "alipay.pay.app.marketing.consult"
}

func (this MarketingConsultParam) Params() map[string]string {
    var m = make(map[string]string)
    //m["biz_scene"] = this.BizScene
    //m["mobile"] = this.Mobile
    //m["encrypted_mobile"] = this.Mobile
    //m["app_auth_token"] = this.AppAuthToken
    return m
}

func (this MarketingConsultParam) ExtJSONParamName() string {
    return "biz_content"
}

func (this MarketingConsultParam) ExtJSONParamValue() string {
    var bytes, err = json.Marshal(this)
    if err != nil {
        return ""
    }
    return string(bytes)
}






