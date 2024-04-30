package translate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// PaPaGoConfigType 啪啪GO翻译配置
type PaPaGoConfigType struct {
	KeyId       string `json:"keyId"`
	Key         string `json:"key"`
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
}

type PaPaGoHTTPTranslateResp struct {
	Message Message `json:"message"`
}

type Message struct {
	Result Result `json:"result"`
}

type Result struct {
	SrcLangType    string `json:"srcLangType"`
	TarLangType    string `json:"tarLangType"`
	TranslatedText string `json:"translatedText"`
}

func (t *PaPaGoConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	if t == nil || t.Key == "" || t.KeyId == "" || t.Url == "" {
		err = errors.New("PaPaGo配置异常")
		return
	}
	mode := t.GetMode()
	// 处理目标语言
	from, err := SafeLangType(req.From, mode)
	if err != nil {
		return
	}
	to, err := SafeLangType(req.To, mode)
	if err != nil {
		return
	}
	// 发起请求
	postResp, err := g.Client().
		SetTimeout(time.Duration(t.CurlTimeOut)*time.Millisecond).
		ContentJson().
		Header(g.MapStrStr{
			"X-NCP-APIGW-API-KEY-ID": t.KeyId,
			"X-NCP-APIGW-API-KEY":    t.Key,
		}).
		Post(context.Background(), t.Url, g.Map{"source": from, "target": to, "text": req.TextStr})
	if err != nil {
		return
	}
	defer func() {
		_ = postResp.Close()
	}()
	bodyByte := postResp.ReadAll()
	// 判断请求状态
	if postResp.StatusCode != 200 {
		err = fmt.Errorf("PaPaGo 请求失败 %d %s", postResp.StatusCode, bodyByte)
		return
	}
	httpResp := new(PaPaGoHTTPTranslateResp)
	if err = json.Unmarshal(bodyByte, httpResp); err != nil {
		return
	}
	lang, err := GetYouDaoLang(httpResp.Message.Result.TarLangType, mode)
	if err != nil {
		return
	}
	for _, v := range strings.Split(httpResp.Message.Result.TranslatedText, "\n") {
		resp = append(resp, &TranslateResp{
			Text:     v,
			FromLang: lang,
		})
	}

	return
}

func (t *PaPaGoConfigType) GetMode() string {
	return PaPaGoTranslateMode
}
