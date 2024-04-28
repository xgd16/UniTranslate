package translate

import (
	"context"
	"errors"
	"fmt"
	"uniTranslate/src/global"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/sashabaranov/go-openai"
)

type ChatGptConfigType struct {
	Key string `json:"key"`
}

func (t *ChatGptConfigType) Translate(from, to, text string) (result []string, fromLang string, err error) {
	if t.Key == "" {
		return nil, "", errors.New("chatGPT翻译配置异常")
	}
	mode := t.GetMode()
	// 语言标记转换
	from, err = SafeLangType(from, mode)
	if err != nil {
		return
	}
	to, err = SafeLangType(to, mode)
	if err != nil {
		return
	}
	// google auto = ""
	if from == "auto" {
		from = ""
	}
	result = make([]string, 0)
	gptResp, err := SendToChatGpt(t.Key, fmt.Sprintf("将[%s]翻译成%s按照格式{\"fromLang\":\"源语言\",\"text\":\"翻译结果\"}返回给我fromLang有这几种语言直接给我返回对应的key位置%s不需要其他任何回复严格按照我给你的格式翻译结果不要用[]包着", text, to, global.ChatGPTLangConfig))
	if err != nil {
		return
	}
	respData := gvar.New(gptResp).MapStrVar()
	result = append(result, respData["text"].String())
	fromLang, err = GetYouDaoLang(respData["fromLang"].String(), mode)
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
