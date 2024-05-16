package translate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
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

type XunFeiNiuConfigType struct {
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
	ApiKey string `json:"apiKey"`
}

// XunFeiNiuTranslate 讯飞新翻译引擎
func (t *XunFeiNiuConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	return xunFeiBaseTranslate(xunFeiNiuHttpConfig, t.GetMode(), &XunFeiConfigType{
		AppId:  t.AppId,
		Secret: t.Secret,
		ApiKey: t.ApiKey,
	}, req.From, req.To, req.Text)
}

func (t *XunFeiNiuConfigType) GetMode() string {
	return XunFeiNiuTranslateMode
}

type XunFeiConfigType struct {
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
	ApiKey string `json:"apiKey"`
}

func (t *XunFeiConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	return xunFeiBaseTranslate(xunFeiHttpConfig, t.GetMode(), t, req.From, req.To, req.Text)
}

func (t *XunFeiConfigType) GetMode() string {
	return XunFeiTranslateMode
}

type XunFeiHTTPTranslateResp struct {
	Code    int64      `json:"code"`
	Data    XunFeiData `json:"data"`
	Message string     `json:"message"`
	Sid     string     `json:"sid"`
}

type XunFeiData struct {
	Result XunFeiResult `json:"result"`
}

type XunFeiResult struct {
	From        string            `json:"from"`
	To          string            `json:"to"`
	TransResult XunFeiTransResult `json:"trans_result"`
}

type XunFeiTransResult struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
}

// xunFeiBaseTranslate 讯飞基础翻译实现
func xunFeiBaseTranslate(baseConfig *xunFeiHttpConfigType, mode string, config *XunFeiConfigType, from, to string, text []string) (resp []*TranslateResp, err error) {
	if config.AppId == "" || config.ApiKey == "" || config.Secret == "" {
		err = errors.New("讯飞翻译配置异常")
		return
	}
	// oFrom := from
	// 语言标记转换
	from, err = SafeLangType(from, mode)
	if err != nil {
		return
	}
	to, err = SafeLangType(to, mode)
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
			"text": base64.StdEncoding.EncodeToString([]byte(strings.Join(text, "№"))),
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
	defer func() { _ = xunFeiResp.Close() }()
	if xunFeiResp.StatusCode != 200 {
		err = fmt.Errorf("请求失败 状态码: %d 返回结果: %s", xunFeiResp.StatusCode, xunFeiResp.ReadAllString())
		return
	}
	respByte := xunFeiResp.ReadAll()
	httpResp := new(XunFeiHTTPTranslateResp)
	if err = json.Unmarshal(respByte, httpResp); err != nil {
		return
	}
	if httpResp.Code != 0 {
		err = fmt.Errorf("讯飞翻译请求失败 %d %s", httpResp.Code, respByte)
		return
	}
	lang, err := GetYouDaoLang(httpResp.Data.Result.From, mode)
	if err != nil {
		return
	}
	for _, item := range strings.Split(httpResp.Data.Result.TransResult.Dst, "№") {
		resp = append(resp, &TranslateResp{
			Text:     item,
			FromLang: lang,
		})
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
