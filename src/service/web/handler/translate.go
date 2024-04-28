package handler

import (
	"uniTranslate/src/translate"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/text/gstr"
)

// Translate 翻译
func Translate(config *types.TranslatePlatform, OriginalFrom, OriginalTo, text string) (data *types.TranslateData, err error) {
	// 获取翻译平台
	t, err := translate.GetTranslate(config.Type, config.Cfg)
	if err != nil {
		return
	}
	// 翻译
	translateTextArr, from, err := t.Translate(OriginalFrom, OriginalTo, text)
	if err != nil {
		return
	}
	// 返回数据
	data = &types.TranslateData{
		OriginalText:     text,
		OriginalTextMd5:  gmd5.MustEncrypt(text),
		TranslateTextArr: translateTextArr,
		From:             from,
		To:               OriginalTo,
		Platform:         config.Type,
		OriginalTextLen:  gstr.LenRune(text),
		TranslationLen:   gstr.LenRune(gstr.Join(translateTextArr, "")),
	}
	return
}
