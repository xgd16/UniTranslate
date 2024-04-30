package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// GoogleConfigType 谷歌配置类型
type GoogleConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	Key         string `json:"key"`
}

type GoogleHTTPTranslateResp struct {
	Data GoogleData `json:"data"`
}

type GoogleData struct {
	Translations []GoogleTranslation `json:"translations"`
}

type GoogleTranslation struct {
	TranslatedText         string `json:"translatedText"`
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
}

// Translate google 翻译
func (t *GoogleConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	if t == nil || t.Url == "" || t.Key == "" {
		err = errors.New("google翻译配置异常")
		return
	}
	mode := t.GetMode()
	ctx := gctx.New()
	// 语言标记转换
	from, err := SafeLangType(req.From, mode)
	if err != nil {
		return
	}
	to, err := SafeLangType(req.To, mode)
	if err != nil {
		return
	}
	if from == "auto" {
		from = ""
	}
	// 发起请求翻译
	postResp, err := g.Client().Discovery(nil).ContentJson().SetTimeout(time.Duration(t.CurlTimeOut)*time.Millisecond).Post(ctx, fmt.Sprintf("%s?key=%s", t.Url, t.Key), g.Map{
		"q":      req.Text,
		"target": to,
		"source": from,
	})
	if err != nil {
		return
	}
	defer func() {
		_ = postResp.Close()
	}()
	respByte := postResp.ReadAll()
	// 处理HTTP状态
	if postResp.StatusCode != 200 {
		err = fmt.Errorf("google翻译请求失败,状态码:%d 返回数据: %s", postResp.StatusCode, respByte)
		return
	}
	// 解析json
	httpResp := new(GoogleHTTPTranslateResp)
	if err = json.Unmarshal(respByte, httpResp); err != nil {
		return
	}
	resp = make([]*TranslateResp, 0)
	for _, item := range httpResp.Data.Translations {
		lang, err1 := GetYouDaoLang(item.DetectedSourceLanguage, mode)
		if err1 != nil {
			err = err1
			return
		}
		resp = append(resp, &TranslateResp{
			Text:     item.TranslatedText,
			FromLang: lang,
		})
	}
	return
}

func (t *GoogleConfigType) GetMode() string {
	return GoogleTranslateMode
}
