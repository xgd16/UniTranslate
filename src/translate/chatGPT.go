package translate

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/sashabaranov/go-openai"
	"github.com/xgd16/gf-x-tool/xtranslate"
	"uniTranslate/src/global"
)

func ChatGptTranslate(config *ChatGptConfigType, from, to, text string) (result []string, fromLang string, err error) {
	if config.Key == "" {
		return nil, "", errors.New("google翻译配置异常")
	}
	// 语言标记转换
	from, err = xtranslate.SafeLangType(from, ChatGptTranslateMode)
	to, err = xtranslate.SafeLangType(to, ChatGptTranslateMode)
	// google auto = ""
	if from == "auto" {
		from = ""
	}
	// 处理转换为安全语言类型错误
	if err != nil {
		return
	}
	// 处理转换后语言设置为auto
	if to == "auto" {
		err = errors.New("转换后语言不能为auto")
		return
	}
	result = make([]string, 0)
	gptResp, err := SendToChatGpt(config.Key, fmt.Sprintf("将[%s]翻译成%s按照格式{\"fromLang\":\"源语言\",\"text\":\"翻译结果\"}返回给我fromLang有这几种语言直接给我返回对应的key位置%s不需要其他任何回复严格按照我给你的格式", text, to, global.ChatGPTLangConfig))
	if err != nil {
		return
	}
	respData := gvar.New(gptResp).MapStrVar()
	result = append(result, respData["text"].String())
	fromLang = respData["fromLang"].String()
	return
}

func SendToChatGpt(key, msg string) (resp string, err error) {
	client := openai.NewClient(key)
	respData, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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
