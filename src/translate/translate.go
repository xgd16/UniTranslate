package translate

import (
	"errors"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gfile"
)

var TranslateModeList = []string{
	BaiduTranslateMode,
	DeeplTranslateMode,
	GoogleTranslateMode,
	YouDaoTranslateMode,
	ChatGptTranslateMode,
	XunFeiTranslateMode,
	XunFeiNiuTranslateMode,
	TencentTranslateMode,
	HuoShanTranslateMode,
	PaPaGoTranslateMode,
}

// ChatGptTranslateMode ChatGPT 支持
const ChatGptTranslateMode = "ChatGPT"

type ChatGptConfigType struct {
	Key string `json:"key"`
}

// XunFeiTranslateMode 讯飞常用版本
const XunFeiTranslateMode = "XunFei"

// XunFeiNiuTranslateMode 讯飞新版
const XunFeiNiuTranslateMode = "XunFeiNiu"

type XunFeiConfigType struct {
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
	ApiKey string `json:"apiKey"`
}

const TencentTranslateMode = "Tencent"

type TencentConfigType struct {
	Url       string `json:"url"`
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
	Region    string `json:"region"`
}

// HuoShanTranslateMode 火山翻译
const HuoShanTranslateMode = "HuoShan"

// HuoShanConfigType 火山翻译配置
type HuoShanConfigType struct {
	AccessKey string
	SecretKey string
}

// PaPaGoTranslateMode 啪啪GO翻译
const PaPaGoTranslateMode = "PaPaGo"

// PaPaGoConfigType 啪啪GO翻译配置
type PaPaGoConfigType struct {
	KeyId       string `json:"keyId"`
	Key         string `json:"key"`
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
}

const (
	YouDaoTranslateMode = "YouDao" // 有道
	BaiduTranslateMode  = "Baidu"  // 百度
	GoogleTranslateMode = "Google" // 谷歌
	DeeplTranslateMode  = "Deepl"  // Deepl
)

// BaiduConfigType 百度的配置类型
type BaiduConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	AppId       string `json:"appId"`
	Key         string `json:"key"`
}

// YouDaoConfigType 有道配置类型
type YouDaoConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	AppKey      string `json:"appKey"`
	SecKey      string `json:"secKey"`
}

// GoogleConfigType 谷歌配置类型
type GoogleConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	Key         string `json:"key"`
}

// DeeplConfigType Deepl配置类型
type DeeplConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	Key         string `json:"key"`
}

// BaseTranslateConf 基础翻译配置
var BaseTranslateConf map[string]map[string]string

// BasePlatformTranslateConf 基础平台翻译配置
var BasePlatformTranslateConf map[string][]map[string]*gvar.Var

// InitTranslateBaseConf 初始化翻译基础配置
var InitTranslateBaseConf = func() (m map[string]map[string]string) {
	// 读取配置文件
	translate := gfile.GetContents("./translate.json")
	if translate == "" {
		return
	}
	// 解析配置文件
	json, err := gjson.DecodeToJson(translate)
	if err != nil {
		return
	}
	// 转换为map
	m = make(map[string]map[string]string, 1)
	for s, v := range json.Var().MapStrVar() {
		m[s] = v.MapStrStr()
	}
	return
}

// InitTranslate 初始化翻译
func InitTranslate() {
	// 初始化基本配置
	BaseTranslateConf = InitTranslateBaseConf()
	if BaseTranslateConf == nil {
		panic("初始化翻译配置失败")
	}
}

// SafeLangType 安全的语言类型
func SafeLangType(t, app string) (lang string, err error) {
	if t == "auto" {
		lang = "auto"
		return
	}
	// 从配置文件中获取
	languages := BaseTranslateConf[app]
	if languages == nil {
		err = errors.New("没有找到应用")
		return
	}
	// 获取语言
	lang = languages[t]
	if lang == "" {
		err = errors.New("不支持的语言类型")
		return
	}
	return
}

// GetYouDaoLang 获取有道语言
func GetYouDaoLang(lang, app string) (youDaoLang string, err error) {
	if lang == "auto" {
		youDaoLang = "auto"
		return
	}
	// 从配置文件中获取
	languages := BaseTranslateConf[app]
	if languages == nil {
		err = errors.New("没有找到应用")
		return
	}
	// 获取语言
	for s, s2 := range languages {
		if s2 == lang {
			youDaoLang = s
			return
		}
	}
	// 获取语言
	youDaoLang = lang
	return
}
