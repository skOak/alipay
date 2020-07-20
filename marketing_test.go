package alipay

import (
    "fmt"
    "testing"
)

//var client *AliPay



func TestAliPay_MarketingConsult(t *testing.T) {
    //client = New(app_id, partner_id, public_key, private_key, false, true)
    client = New(this_appID, this_partnerID, this_publicKey, this_privateKey, false, true)
    client.AliPayPublicKey = this_aliPublicKey

    fmt.Println("========== MarketingConsult ==========")
    var p = MarketingConsultParam{}
    p.BizScene = "ORDER_PAGE"
    p.Mobile = "13333333333"
    p.EncryptedMobile = "34347C343003E57232A5D21F14FE399E"
    p.TotalAmount = "88"

    fmt.Println(client.MarketingConsult(p))
}