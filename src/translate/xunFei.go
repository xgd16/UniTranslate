package translate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/xgd16/gf-x-tool/xtranslate"
)

type xunFeiHttpConfigType struct {
	Url  string
	Host string
	Uri  string
}

var xunFeiHttpConfig = &xunFeiHttpConfigType{}
var xunFeiNiuHttpConfig = &xunFeiHttpConfigType{}

func init() {
	xunFeiHttpConfig.Url = "https://itrans.xfyun.cn/v2/its"
	parse, _ := url.Parse(xunFeiHttpConfig.Url)
	xunFeiHttpConfig.Host = parse.Host
	xunFeiHttpConfig.Uri = parse.Path

	xunFeiNiuHttpConfig.Url = "https://ntrans.xfyun.cn/v2/ots"
	parse, _ = url.Parse(xunFeiNiuHttpConfig.Url)
	xunFeiNiuHttpConfig.Host = parse.Host
	xunFeiNiuHttpConfig.Uri = parse.Path
}

// XunFeiNiuTranslate 讯飞新翻译引擎
func XunFeiNiuTranslate(config *XunFeiConfigType, from, to, text string) (result []string, fromLang string, err error) {
	return xunFeiBaseTranslate(xunFeiNiuHttpConfig, XunFeiNiuTranslateMode, config, from, to, text)
}

func XunFeiTranslate(config *XunFeiConfigType, from, to, text string) (result []string, fromLang string, err error) {
	return xunFeiBaseTranslate(xunFeiHttpConfig, XunFeiTranslateMode, config, from, to, text)
}

// xunFeiBaseTranslate 讯飞基础翻译实现
func xunFeiBaseTranslate(baseConfig *xunFeiHttpConfigType, mode string, config *XunFeiConfigType, from, to, text string) (result []string, fromLang string, err error) {
	if config.AppId == "" || config.ApiKey == "" || config.Secret == "" {
		return nil, "", errors.New("讯飞翻译配置异常")
	}
	oFrom := from
	// 语言标记转换
	from, err = xtranslate.SafeLangType(from, mode)
	if err != nil {
		return
	}
	to, err = xtranslate.SafeLangType(to, mode)
	if err != nil {
		return
	}
	if mode == XunFeiTranslateMode && from == "auto" {
		err = errors.New("当前翻译平台不支持自动识别源语言语种")
		return
	}
	// http client
	postData := g.Map{
		"common": map[string]interface{}{
			"app_id": config.AppId,
		},
		"business": map[string]interface{}{
			"from": from,
			"to":   to,
		},
		"data": map[string]interface{}{
			"text": base64.StdEncoding.EncodeToString([]byte(text)),
		},
	}
	currentTime := time.Now().UTC().Format(time.RFC1123)
	// get Digest
	tt := gjson.MustEncodeString(postData)
	digest := fmt.Sprintf("SHA-256=%s", xunFeiSignBody(tt))
	// get sign
	sign := xunFeiGenerateSignature(baseConfig.Host, currentTime, "POST", baseConfig.Uri, "HTTP/1.1", digest, config.Secret)
	// send request
	xunFeiResp, err := gclient.New().ContentJson().Header(g.MapStrStr{
		"Date":          currentTime,
		"Digest":        digest,
		"Authorization": fmt.Sprintf(`api_key="%s", algorithm="%s", headers="host date request-line digest", signature="%s"`, config.ApiKey, "hmac-sha256", sign),
	}).Post(gctx.New(), baseConfig.Url, postData)
	if err != nil {
		return
	}
	// to json
	jsonData, err := gjson.DecodeToJson(xunFeiResp.ReadAllString())
	if err != nil {
		return
	}
	if jsonData.Get("code").Int() != 0 {
		err = errors.New(jsonData.Get("message").String())
		return
	}
	result = []string{gstr.Trim(jsonData.Get("data.result.trans_result.dst").String())}
	fromLang, err = xtranslate.GetYouDaoLang(jsonData.Get("data.result.from", "").String(), XunFeiNiuTranslateMode)
	if fromLang == "" {
		fromLang = oFrom
	}
	return
}

func xunFeiGenerateSignature(host, date, httpMethod, requestUri, httpProto, digest string, secret string) string {
	var signatureStr string
	if len(host) != 0 {
		signatureStr = "host: " + host + "\n"
	}
	signatureStr += "date: " + date + "\n"
	signatureStr += httpMethod + " " + requestUri + " " + httpProto + "\n"
	signatureStr += "digest: " + digest
	return xunFeiHmacsign(signatureStr, secret)
}

func xunFeiHmacsign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func xunFeiSignBody(data string) string {
	// 进行sha256签名
	sha := sha256.New()
	sha.Write([]byte(data))
	encodeData := sha.Sum(nil)
	// 经过base64转换
	return base64.StdEncoding.EncodeToString(encodeData)
}
