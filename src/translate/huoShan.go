package translate

import (
	"encoding/json"
	baseErr "errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

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

// HuoShanConfigType 火山翻译配置
type HuoShanConfigType struct {
	AccessKey string
	SecretKey string
}

type HuoShanHTTPTranslateResp struct {
	TranslationList  []TranslationList `json:"TranslationList"`
	ResponseMetadata ResponseMetadata  `json:"ResponseMetadata"`
	ResponseMetaData ResponseMetadata  `json:"ResponseMetaData"`
}

type ResponseMetadata struct {
	RequestID string `json:"RequestId"`
	Action    string `json:"Action"`
	Version   string `json:"Version"`
	Service   string `json:"Service"`
	Region    string `json:"Region"`
}

type TranslationList struct {
	Translation            string      `json:"Translation"`
	DetectedSourceLanguage string      `json:"DetectedSourceLanguage"`
	Extra                  interface{} `json:"Extra"`
}

// Translate 火山翻译
func (t *HuoShanConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	if t == nil || t.SecretKey == "" || t.AccessKey == "" {
		err = baseErr.New("火山翻译配置异常")
		return
	}
	mode := t.GetMode()
	from := req.From
	// 处理目标语言
	if from == "auto" {
		from = ""
	} else {
		from, err = SafeLangType(from, mode)
		if err != nil {
			return
		}
	}
	to, err := SafeLangType(req.To, mode)
	if err != nil {
		return
	}

	client := base.NewClient(ServiceInfo, ApiInfoList)
	client.SetAccessKey(t.AccessKey)
	client.SetSecretKey(t.SecretKey)

	body, err := json.Marshal(Req{
		SourceLanguage: from,
		TargetLanguage: to,
		TextList:       req.Text,
	})
	if err != nil {
		return
	}
	respByte, code, err := client.Json("TranslateText", nil, string(body))
	if err != nil {
		return
	}
	if code != 200 {
		err = fmt.Errorf("火山翻译请求失败,状态码:%d 返回数据: %s", code, respByte)
		return
	}
	httpResp := new(HuoShanHTTPTranslateResp)
	if err = json.Unmarshal(respByte, httpResp); err != nil {
		return
	}
	for _, item := range httpResp.TranslationList {
		lang, err1 := GetYouDaoLang(item.DetectedSourceLanguage, mode)
		if err1 != nil {
			err = err1
			return
		}
		resp = append(resp, &TranslateResp{
			Text:     item.Translation,
			FromLang: lang,
		})

	}
	return
}

func (t *HuoShanConfigType) GetMode() string {
	return HuoShanTranslateMode
}
