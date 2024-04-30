package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// DeeplConfigType Deepl配置类型
type DeeplConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	Key         string `json:"key"`
}

type DeeplHTTPTranslateResp struct {
	Translations []DeeplTranslation `json:"translations"`
}

type DeeplTranslation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

func (t *DeeplConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	if t == nil || t.Url == "" || t.Key == "" {
		err = errors.New("deepl翻译配置异常")
		return
	}
	ctx := gctx.New()
	mode := t.GetMode()
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
	// 调用翻译
	HttpResult, err := g.Client().SetTimeout(time.Duration(t.CurlTimeOut)*time.Millisecond).Header(g.MapStrStr{
		"Authorization": fmt.Sprintf("DeepL-Auth-Key %s", t.Key),
	}).Post(ctx, t.Url, g.Map{
		"text":        req.Text,
		"source_lang": from,
		"target_lang": to,
	})
	// 处理调用接口错误
	if err != nil {
		return
	}
	// 推出函数时关闭链接
	defer func() { _ = HttpResult.Close() }()
	// 判断状态码
	respByte := HttpResult.ReadAll()
	if HttpResult.StatusCode != 200 {
		err = fmt.Errorf("请求失败 状态码: %d 返回结果: %s", HttpResult.StatusCode, respByte)
		return
	}
	fmt.Printf("deepl resp: %s\n", respByte)
	httpResp := new(DeeplHTTPTranslateResp)
	if err = json.Unmarshal(respByte, httpResp); err != nil {
		return
	}
	if len(httpResp.Translations) == 0 {
		return
	}
	respData := httpResp.Translations[0]

	lang, err := GetYouDaoLang(respData.DetectedSourceLanguage, mode)
	if err != nil {
		return
	}
	strArr := make([]string, 0)
	if err = json.Unmarshal([]byte(respData.Text), &strArr); err != nil {
		return
	}
	for _, v := range strArr {
		resp = append(resp, &TranslateResp{
			Text:     v,
			FromLang: lang,
		})
	}
	return
}

func (t *DeeplConfigType) GetMode() string {
	return DeeplTranslateMode
}
