package translate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

var ChatGPTLangConfig string

type ChatGptConfigType struct {
	Key string `json:"key"`
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
	gptResp, err := SendToChatGpt(t.Key, fmt.Sprintf("接下来模拟你是一个批量翻译接口你将 [%s] 数组数据批量翻译为 %s 返回数据结构为[ [{\"fromLang\":\"源语言\",\"text\":\"翻译结果\"}] ]标准的压缩数组json结构返回给我fromLang有这几种语言直接给我返回对应的key位置%s不需要其他任何回复严格按照我给你的格式翻译结果不要用[]包着", func() (str string) {
		for k, v := range req.Text {
			if k == 0 {
				str = fmt.Sprintf("\"%s\"", v)
			} else {
				str = fmt.Sprintf("%s,\"%s\"", str, v)
			}
		}
		return
	}(), to, ChatGPTLangConfig))
	if err != nil {
		return
	}
	httpResp := new(ChatGPTHTTPTranslateResp)
	if err = json.Unmarshal([]byte(gptResp), httpResp); err != nil {
		return
	}
	resp = make([]*TranslateResp, 0)
	for _, item := range *httpResp {
		resp = append(resp, &TranslateResp{
			Text:     item.Text,
			FromLang: item.FromLang,
		})
	}
	return
}

func SendToChatGpt(key, msg string) (resp string, err error) {
	client := openai.NewClient(key)
	respData, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)
	if err != nil {
		return
	}
	resp = respData.Choices[0].Message.Content
	return
}

func (t *ChatGptConfigType) GetMode() string {
	return ChatGptTranslateMode
}
