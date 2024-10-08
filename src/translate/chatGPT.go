package translate

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

var ChatGPTLangConfig string

type ChatGptConfigType struct {
	Key   string `json:"key"`
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
		respStr, err1 := SendToChatGpt(t.Key, t.OrgId, fmt.Sprintf("Translate the following text to %s and output only the translation: %s", to, item), t.Model)
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

func SendToChatGpt(key, orgId, msg, modelStr string) (resp string, err error) {
	config := openai.DefaultConfig(key)
	if orgId != "" {
		config.OrgID = orgId
	}
	client := openai.NewClientWithConfig(config)
	model := openai.GPT3Dot5Turbo0125
	switch modelStr {
	case "gpt-3.5-turbo-0125":
		model = openai.GPT3Dot5Turbo0125
	case "gpt-4-turbo":
		model = openai.GPT4Turbo
	case "gpt-3.5-turbo":
		model = openai.GPT3Dot5Turbo
	}
	respData, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
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
