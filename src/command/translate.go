package command

import (
	"context"
	"fmt"
	"uniTranslate/src/logic"
	"uniTranslate/src/types"

	"github.com/fatih/color"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/xgd16/gf-x-tool/xlib"
)

func Translate(ctx context.Context, parser *gcmd.Parser) (err error) {
	from := parser.GetArg(2).String()
	to := parser.GetArg(3).String()
	if from == "" || to == "" {
		fmt.Println("参数错误")
		return
	}
	// 开启交互
	for {
		text := gcmd.Scan(color.Set(color.FgHiGreen).Sprintf("请输入要翻译的内容: "))
		if text == "" {
			color.Set(color.FgHiYellow).Println("请输入内容!")
			continue
		}
		data, err := logic.Translate(ctx, "终端", &types.TranslateReq{
			From: from,
			To:   to,
			Text: []string{text},
		})
		if err != nil || len(data.Translate) <= 0 {
			color.Set(color.FgHiRed).Println("翻译失败请重试")
			continue
		}
		t := data.Translate[0]
		fmt.Printf("翻译结果: %s\n语言类型: %s -> %s\n", t.Text, xlib.IF(from == "auto", t.FromLang, from), to)
	}
}