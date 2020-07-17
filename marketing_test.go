package alipay

import (
    "fmt"
    "testing"
)

//var client *AliPay

func TestAliPay_MarketingConsult(t *testing.T) {
    client = New(this_appID, this_partnerID, this_publicKey, this_privateKey, false, true)
    client.AliPayPublicKey = this_aliPublicKey

    fmt.Println("========== MarketingConsult ==========")
    var p = MarketingConsultParam{}
    p.BizScene = "DETAIL_PAGE"

    fmt.Println(client.MarketingConsult(p))
}