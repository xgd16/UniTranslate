package handler

import (
	"errors"
	"uniTranslate/src/translate"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func Translate(config *types.TranslatePlatform, OriginalFrom, OriginalTo, text string) (data *types.TranslateData, translateErr error) {
	var (
		translateTextArr []string
		from             string
	)
	// 调用对应平台
	switch config.Type {
	case translate.BaiduTranslateMode:
		translateTextArr, from, translateErr = translate.BaiduTranslate(&translate.BaiduConfigType{
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
			AppId:       gconv.String(config.Cfg["appId"]),
			Key:         gconv.String(config.Cfg["key"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.YouDaoTranslateMode:
		translateTextArr, from, translateErr = translate.YouDaoTranslate(&translate.YouDaoConfigType{
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
			AppKey:      gconv.String(config.Cfg["appKey"]),
			SecKey:      gconv.String(config.Cfg["secKey"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.GoogleTranslateMode:
		translateTextArr, from, translateErr = translate.GoogleTranslate(&translate.GoogleConfigType{
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
			Key:         gconv.String(config.Cfg["key"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.DeeplTranslateMode:
		translateTextArr, from, translateErr = translate.DeeplTranslate(&translate.DeeplConfigType{
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
			Key:         gconv.String(config.Cfg["key"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.ChatGptTranslateMode:
		translateTextArr, from, translateErr = translate.ChatGptTranslate(&translate.ChatGptConfigType{
			Key: gconv.String(config.Cfg["key"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.XunFeiTranslateMode:
		translateTextArr, from, translateErr = translate.XunFeiTranslate(&translate.XunFeiConfigType{
			AppId:  gconv.String(config.Cfg["appId"]),
			ApiKey: gconv.String(config.Cfg["apiKey"]),
			Secret: gconv.String(config.Cfg["secret"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.XunFeiNiuTranslateMode:
		translateTextArr, from, translateErr = translate.XunFeiNiuTranslate(&translate.XunFeiConfigType{
			AppId:  gconv.String(config.Cfg["appId"]),
			ApiKey: gconv.String(config.Cfg["apiKey"]),
			Secret: gconv.String(config.Cfg["secret"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.TencentTranslateMode:
		translateTextArr, from, translateErr = translate.TencentTranslate(&translate.TencentConfigType{
			Url:       gconv.String(config.Cfg["url"]),
			SecretId:  gconv.String(config.Cfg["secretId"]),
			SecretKey: gconv.String(config.Cfg["secretKey"]),
			Region:    gconv.String(config.Cfg["region"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.HuoShanTranslateMode:
		translateTextArr, from, translateErr = translate.HuoShanTranslate(&translate.HuoShanConfigType{
			AccessKey: gconv.String(config.Cfg["accessKey"]),
			SecretKey: gconv.String(config.Cfg["secretKey"]),
		}, OriginalFrom, OriginalTo, text)
	case translate.PaPaGoTranslateMode:
		translateTextArr, from, translateErr = translate.PaPaGoTranslate(&translate.PaPaGoConfigType{
			KeyId:       gconv.String(config.Cfg["keyId"]),
			Key:         gconv.String(config.Cfg["key"]),
			CurlTimeOut: gconv.Int(config.Cfg["curlTimeOut"]),
			Url:         gconv.String(config.Cfg["url"]),
		}, OriginalFrom, OriginalTo, text)
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
