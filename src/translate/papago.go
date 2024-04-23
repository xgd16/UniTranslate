package translate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/xgd16/gf-x-tool/xtranslate"
)

func PaPaGoTranslate(config *PaPaGoConfigType, from, to, text string) (result []string, fromLang string, err error) {
	if config == nil || config.Key == "" || config.KeyId == "" || config.Url == "" {
		err = errors.New("PaPaGo配置异常")
		return
	}
	// 处理目标语言
	from, err = xtranslate.SafeLangType(from, PaPaGoTranslateMode)
	if err != nil {
		return
	}
	to, err = xtranslate.SafeLangType(to, PaPaGoTranslateMode)
	if err != nil {
		return
	}
	// 发起请求
	resp, err := g.Client().
		SetTimeout(time.Duration(config.CurlTimeOut)*time.Millisecond).
		ContentJson().
		Header(g.MapStrStr{
			"X-NCP-APIGW-API-KEY-ID": config.KeyId,
			"X-NCP-APIGW-API-KEY":    config.Key,
		}).
		Post(context.Background(), config.Url, g.Map{"source": from, "target": to, "text": text})
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
	fromLang, err = xtranslate.GetYouDaoLang(jsonData.Get("message.result.srcLangType").String(), PaPaGoTranslateMode)
	return
}
