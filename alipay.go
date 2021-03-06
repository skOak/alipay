package alipay

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/skOak/alipay/encoding"
)

type AliPay struct {
	appId           string
	apiDomain       string
	partnerId       string
	publicKey       []byte
	privateKey      []byte
	AliPayPublicKey []byte
	client          *http.Client
	SignType        string
}

func New(appId, partnerId string, publicKey, privateKey []byte, isProduction, isMApi bool) (client *AliPay) {
	client = &AliPay{}
	client.appId = appId
	client.partnerId = partnerId
	client.privateKey = privateKey
	client.publicKey = publicKey
	client.client = http.DefaultClient
	if isMApi {
		client.apiDomain = K_MAPI_PRODUCTION_API_URL
		client.SignType = K_SIGN_TYPE_RSA
	} else {
		if isProduction {
			client.apiDomain = K_ALI_PAY_PRODUCTION_API_URL
		} else {
			client.apiDomain = K_ALI_PAY_SANDBOX_API_URL
		}
		client.SignType = K_SIGN_TYPE_RSA2
	}
	return client
}

func (this *AliPay) URLValues(param AliPayParam) (value url.Values, err error) {
	var p = url.Values{}
	p.Add("app_id", this.appId)
	p.Add("method", param.APIName())
	p.Add("format", K_FORMAT)
	p.Add("charset", K_CHARSET)
	p.Add("sign_type", this.SignType)
	p.Add("timestamp", time.Now().Format(K_TIME_FORMAT))
	p.Add("version", K_VERSION)

	if len(param.ExtJSONParamName()) > 0 {
		p.Add(param.ExtJSONParamName(), param.ExtJSONParamValue())
	}

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
	p.Add("sign", sign)

	if err != nil {
		return nil, err
	}
	return p, nil
}

type HasOriginBodySet interface {
	SetOriginBody(body string)
}

func (this *AliPay) doRequest(method string, param AliPayParam, results interface{}) (err error) {
	var buf io.Reader
	if param != nil {
		p, err := this.URLValues(param)
		if err != nil {
			return err
		}
		buf = strings.NewReader(p.Encode())
	}

	req, err := http.NewRequest(method, this.apiDomain, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := this.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}


	if len(this.AliPayPublicKey) > 0 {
		var dataStr = string(data)

		var rootNodeName = strings.Replace(param.APIName(), ".", "_", -1) + k_RESPONSE_SUFFIX

		var rootIndex = strings.LastIndex(dataStr, rootNodeName)
		var errorIndex = strings.LastIndex(dataStr, k_ERROR_RESPONSE)

		var content string
		var sign string

		if rootIndex > 0 {
			content, sign = parserJSONSource(dataStr, rootNodeName, rootIndex)
		} else if errorIndex > 0 {
			content, sign = parserJSONSource(dataStr, k_ERROR_RESPONSE, errorIndex)
		} else {
			return errors.New("empty response data")
		}

		if ok, err := verifyResponseData([]byte(content), this.SignType, sign, this.AliPayPublicKey); ok == false {
			return err
		}
	}

	err = json.Unmarshal(data, results)
	if err != nil {
		return err
	}
	if t, ok := results.(HasOriginBodySet); ok {
		t.SetOriginBody(string(data))
	}

	return err
}

func (this *AliPay) DoRequest(method string, param AliPayParam, results interface{}) (err error) {
	return this.doRequest(method, param, results)
}

func parserJSONSource(rawData string, nodeName string, nodeIndex int) (content string, sign string) {
	var dataStartIndex = nodeIndex + len(nodeName) + 2
	var signIndex = strings.LastIndex(rawData, "\""+k_SIGN_NODE_NAME+"\"")
	var dataEndIndex = signIndex - 1

	var indexLen = dataEndIndex - dataStartIndex
	if indexLen < 0 {
		return "", ""
	}
	content = rawData[dataStartIndex:dataEndIndex]

	var signStartIndex = signIndex + len(k_SIGN_NODE_NAME) + 4
	sign = rawData[signStartIndex:]
	var signEndIndex = strings.LastIndex(sign, "\"}")
	sign = sign[:signEndIndex]

	return content, sign
}

func signRSA2(keys []string, param url.Values, privateKey []byte) (s string, err error) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	sig, err := encoding.SignPKCS1v15([]byte(src), privateKey, crypto.SHA256)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

func signRSA(keys []string, param url.Values, privateKey []byte) (s string, err error) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	sig, err := encoding.SignPKCS1v15([]byte(src), privateKey, crypto.SHA1)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

func verifySign(req *http.Request, key []byte) (ok bool, err error) {
	sign, err := base64.StdEncoding.DecodeString(req.Form.Get("sign"))
	signType := req.Form.Get("sign_type")
	if err != nil {
		return false, err
	}

	var keys = make([]string, 0, 0)
	for key, value := range req.Form {
		if key == "sign" || key == "sign_type" {
			continue
		}
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		//var value = strings.TrimSpace(req.Form.Get(key))
		// 不需要去除参数中的空格，只需要判断参数是否空值(空字符串)即可
		var value = req.Form.Get(key)
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var s = strings.Join(pList, "&")

	if signType == K_SIGN_TYPE_RSA {
		err = encoding.VerifyPKCS1v15([]byte(s), sign, key, crypto.SHA1)
	} else {
		err = encoding.VerifyPKCS1v15([]byte(s), sign, key, crypto.SHA256)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func verifyResponseData(data []byte, signType, sign string, key []byte) (ok bool, err error) {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	if signType == K_SIGN_TYPE_RSA {
		err = encoding.VerifyPKCS1v15(data, signBytes, key, crypto.SHA1)
	} else {
		err = encoding.VerifyPKCS1v15(data, signBytes, key, crypto.SHA256)
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
