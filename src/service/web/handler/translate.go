package handler

import (
	"errors"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/xgd16/gf-x-tool/xtranslate"
	"uniTranslate/src/translate"
	"uniTranslate/src/types"
)

func Translate(config *types.TranslatePlatform, OriginalFrom, OriginalTo, text string) (data *types.TranslateData, translateErr error) {
	var (
		translateTextArr []string
		from             string
	)
	// 调用对应平台
	switch config.Type {
	case xtranslate.Baidu:
		translateTextArr, from, translateErr = xtranslate.BaiduTranslate(&xtranslate.BaiduConfigType{
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
			AppId:       gconv.String(config.Cfg["appId"]),
			Key:         gconv.String(config.Cfg["key"]),
		}, OriginalFrom, OriginalTo, text)
		break
	case xtranslate.YouDao:
		translateTextArr, from, translateErr = xtranslate.YouDaoTranslate(&xtranslate.YouDaoConfigType{
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
			AppKey:      gconv.String(config.Cfg["appKey"]),
			SecKey:      gconv.String(config.Cfg["secKey"]),
		}, OriginalFrom, OriginalTo, text)
		break
	case xtranslate.Google:
		translateTextArr, from, translateErr = xtranslate.GoogleTranslate(&xtranslate.GoogleConfigType{
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
			Key:         gconv.String(config.Cfg["key"]),
		}, OriginalFrom, OriginalTo, text)
	case xtranslate.Deepl:
		translateTextArr, from, translateErr = xtranslate.DeeplTranslate(&xtranslate.DeeplConfigType{
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
			Key:         gconv.String(config.Cfg["key"]),
		}, OriginalFrom, OriginalTo, text)
		break
	case translate.ChatGptTranslateMode:
		translateTextArr, from, translateErr = translate.ChatGptTranslate(&translate.ChatGptConfigType{
			Key: gconv.String(config.Cfg["key"]),
		}, OriginalFrom, OriginalTo, text)
		break
	case translate.XunFeiTranslateMode:
		translateTextArr, from, translateErr = translate.XunFeiTranslate(&translate.XunFeiConfigType{
			AppId:  gconv.String(config.Cfg["appId"]),
			ApiKey: gconv.String(config.Cfg["apiKey"]),
			Secret: gconv.String(config.Cfg["secret"]),
		}, OriginalFrom, OriginalTo, text)
		break
	case translate.XunFeiNiuTranslateMode:
		translateTextArr, from, translateErr = translate.XunFeiNiuTranslate(&translate.XunFeiConfigType{
			AppId:  gconv.String(config.Cfg["appId"]),
			ApiKey: gconv.String(config.Cfg["apiKey"]),
			Secret: gconv.String(config.Cfg["secret"]),
		}, OriginalFrom, OriginalTo, text)
		break
	case translate.TencentTranslateMode:
		translateTextArr, from, translateErr = translate.TencentTranslate(&translate.TencentConfigType{
			Url:       gconv.String(config.Cfg["url"]),
			SecretId:  gconv.String(config.Cfg["secretId"]),
			SecretKey: gconv.String(config.Cfg["secretKey"]),
			Region:    gconv.String(config.Cfg["region"]),
		}, OriginalFrom, OriginalTo, text)
		break
	case translate.HuoShanTranslateMode:
		translateTextArr, from, translateErr = translate.HuoShanTranslate(&translate.HuoShanConfigType{
			AccessKey: gconv.String(config.Cfg["accessKey"]),
			SecretKey: gconv.String(config.Cfg["secretKey"]),
		}, OriginalFrom, OriginalTo, text)
		break
	default:
		translateErr = errors.New("不支持的翻译")
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
