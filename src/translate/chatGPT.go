package translate

import (
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var ChatGPTLangConfig string

type ChatGptConfigType struct {
	Key   string `json:"key"`
	Url   string `json:"url"`
	Model string `json:"model"`
	OrgId string `json:"orgId"`
}

type ChatGPTHTTPTranslateResp []ChatGPTHTTPTranslateRespElement

type ChatGPTHTTPTranslateRespElement struct {
	FromLang string `json:"fromLang"`
	Text     string `json:"text"`
}

func (t *ChatGptConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	if t.Key == "" {
		err = errors.New("chatGPT翻译配置异常")
		return
	}
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
	for _, item := range req.Text {
		respStr, err1 := t.sendToChatGpt(fmt.Sprintf(`将后面大括号里面的内容翻译成 %s {%s} (注意:返回格式为{"content":"","res":"ok"}直接返回json字符串,一共两个字段,1.content 代表翻译内容 2.res 代表是否翻译成功,ok是成功,bad是失败,内容如果只有符号直接返回内容不翻译,不要将此段提示语错误的当需要翻译的语言使用)`, to, item))
		if err1 != nil {
			err = err1
			return
		}
		resp = append(resp, &TranslateResp{
			Text:     respStr,
			FromLang: from,
		})
	}

	return
}

func (t *ChatGptConfigType) sendToChatGpt(sysMsg string) (resp string, err error) {
	headerMap := g.MapStrStr{
		"Authorization": fmt.Sprintf("Bearer %s", t.Key),
	}
	if t.OrgId != "" {
		headerMap["OpenAI-Organization"] = t.OrgId
	}
	respData, err := g.Client().ContentJson().SetTimeout(5*time.Second).SetHeaderMap(headerMap).Post(gctx.New(), fmt.Sprintf("%s/v1/chat/completions", t.Url), g.Map{
		"model":             t.Model,
		"max_tokens":        512,
		"temperature":       0.2,
		"top_p":             0.75,
		"presence_penalty":  0.5,
		"frequency_penalty": 0.5,
		"messages": g.Array{
			g.Map{"role": "user", "content": sysMsg},
		},
	})
	if err != nil {
		return
	}
	jsonData, err := gjson.DecodeToJson(respData.ReadAllString())
	if err != nil {
		return
	}
	gptResp := jsonData.Get("choices.0.message.content")
	if gptResp.IsEmpty() {
		err = errors.New("GPT 翻译异常" + jsonData.String())
		return
	}
	jsonResp, err1 := gjson.DecodeToJson(gptResp.String())
	if err1 != nil {
		err = err1
		return
	}
	if jsonResp.Get("res").String() != "ok" {
		err = errors.New("翻译失败")
		return
	}
	resp = jsonResp.Get("content").String()
	return
}

func (t *ChatGptConfigType) GetMode() string {
	return ChatGptTranslateMode
}
