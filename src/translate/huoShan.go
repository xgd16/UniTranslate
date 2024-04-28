package translate

import (
	"encoding/json"
	baseErr "errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/pkg/errors"
	"github.com/volcengine/volc-sdk-golang/base"
)

var (
	ServiceInfo = &base.ServiceInfo{
		Timeout: 5 * time.Second,
		Host:    "translate.volcengineapi.com",
		Header: http.Header{
			"Accept": []string{"application/json"},
		},
		Credentials: base.Credentials{Region: base.RegionCnNorth1, Service: "translate"},
	}
	ApiInfoList = map[string]*base.ApiInfo{
		"TranslateText": {
			Method: http.MethodPost,
			Path:   "/",
			Query: url.Values{
				"Action":  []string{"TranslateText"},
				"Version": []string{"2020-06-01"},
			},
		},
	}
)

type Req struct {
	SourceLanguage string   `json:"SourceLanguage"`
	TargetLanguage string   `json:"TargetLanguage"`
	TextList       []string `json:"TextList"`
}

// HuoShanTranslate 火山翻译
func HuoShanTranslate(config *HuoShanConfigType, from, to, text string) (result []string, fromLang string, err error) {
	if config.SecretKey == "" || config.AccessKey == "" {
		err = baseErr.New("火山翻译配置异常")
		return
	}
	// 处理目标语言
	if from == "auto" {
		from = ""
	} else {
		from, err = SafeLangType(from, HuoShanTranslateMode)
		if err != nil {
			return
		}
	}
	to, err = SafeLangType(to, HuoShanTranslateMode)
	if err != nil {
		return
	}

	client := base.NewClient(ServiceInfo, ApiInfoList)
	client.SetAccessKey(config.AccessKey)
	client.SetSecretKey(config.SecretKey)

	req := Req{
		SourceLanguage: from,
		TargetLanguage: to,
		TextList:       []string{text},
	}
	body, err := json.Marshal(req)
	if err != nil {
		return
	}
	resp, code, err := client.Json("TranslateText", nil, string(body))
	if err != nil {
		return
	}
	if code != 200 {
		err = errors.New(string(resp))
		return
	}
	jsonData, err := gjson.DecodeToJson(resp)
	if err != nil {
		return
	}
	result = jsonData.Get("TranslationList.0.Translation").Strings()
	fromLang, err = GetYouDaoLang(jsonData.Get("TranslationList.0.DetectedSourceLanguage").String(), HuoShanTranslateMode)
	return
}
