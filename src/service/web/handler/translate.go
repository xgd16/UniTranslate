package handler

import (
	"uniTranslate/src/translate"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/text/gstr"
)

// Translate 翻译
func Translate(config *types.TranslatePlatform, req *translate.TranslateReq) (data *types.TranslateData, err error) {
	// 获取翻译平台
	t, err := translate.GetTranslate(config.Type, config.Cfg)
	if err != nil {
		return
	}
	// 翻译
	resp, err := t.Translate(req)
	if err != nil {
		return
	}
	// 返回数据
	data = &types.TranslateData{
		OriginalText:    req.Text,
		OriginalTextStr: &req.TextStr,
		OriginalTextMd5: gmd5.MustEncrypt(req.TextStr),
		Translate:       resp,
		To:              req.To,
		Platform:        config.Type,
		OriginalTextLen: gstr.LenRune(gstr.Replace(req.TextStr, "\n", "")),
	}
	return
}
