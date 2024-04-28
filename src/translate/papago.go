package translate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// PaPaGoConfigType 啪啪GO翻译配置
type PaPaGoConfigType struct {
	KeyId       string `json:"keyId"`
	Key         string `json:"key"`
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
}

func (t *PaPaGoConfigType) Translate(from, to, text string) (result []string, fromLang string, err error) {
	if t == nil || t.Key == "" || t.KeyId == "" || t.Url == "" {
		err = errors.New("PaPaGo配置异常")
		return
	}
	mode := t.GetMode()
	// 处理目标语言
	from, err = SafeLangType(from, mode)
	if err != nil {
		return
	}
	to, err = SafeLangType(to, mode)
	if err != nil {
		return
	}
	// 发起请求
	resp, err := g.Client().
		SetTimeout(time.Duration(t.CurlTimeOut)*time.Millisecond).
		ContentJson().
		Header(g.MapStrStr{
			"X-NCP-APIGW-API-KEY-ID": t.KeyId,
			"X-NCP-APIGW-API-KEY":    t.Key,
		}).
		Post(context.Background(), t.Url, g.Map{"source": from, "target": to, "text": text})
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Close()
	}()
	bodyStr := resp.ReadAllString()
	// 判断请求状态
	if resp.StatusCode != 200 {
		err = fmt.Errorf("PaPaGo 请求失败 %d %s", resp.StatusCode, bodyStr)
		return
	}
	// 转换json
	jsonData, err := gjson.DecodeToJson(bodyStr)
	if err != nil {
		return
	}
	if !jsonData.Get("error").IsEmpty() {
		err = errors.New(jsonData.Get("error.message").String())
		return
	}
	respTextT := jsonData.Get("message.result.translatedText")
	if respTextT.IsEmpty() {
		err = fmt.Errorf("PaPaGo 返回结果错误 %s", bodyStr)
		return
	}
	result = []string{respTextT.String()}
	fromLang, err = GetYouDaoLang(jsonData.Get("message.result.srcLangType").String(), mode)
	return
}

func (t *PaPaGoConfigType) GetMode() string {
	return PaPaGoTranslateMode
}
