package translate

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/util/gconv"

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
	FreeGoogleMode,
}

// ITranslate 翻译接口
type ITranslate interface {
	// Translate 翻译
	Translate(req *TranslateReq) (resp []*TranslateResp, err error)
	// GetMode 获取模式
	GetMode() (mode string)
}

type TranslateResp struct {
	Text     string `json:"text"`
	FromLang string `json:"fromLang"`
}

type TranslateHttpReq struct {
	ClientIp string          `json:"clientIp"`
	Context  context.Context `json:"context"`
}

type TranslateReq struct {
	HttpReq  *TranslateHttpReq
	From     string   `json:"from"`
	To       string   `json:"to"`
	Platfrom string   `json:"platform"`
	Text     []string `json:"text"`
	TextStr  string   `json:"textStr"`
}

func GetTranslate(mode string, config map[string]any) (t ITranslate, err error) {
	// 调用对应平台
	switch mode {
	case BaiduTranslateMode:
		t = new(BaiduConfigType)
	case YouDaoTranslateMode:
		t = new(YouDaoConfigType)
	case GoogleTranslateMode:
		t = new(GoogleConfigType)
	case DeeplTranslateMode:
		t = new(DeeplConfigType)
	case ChatGptTranslateMode:
		t = new(ChatGptConfigType)
	case XunFeiTranslateMode:
		t = new(XunFeiConfigType)
	case XunFeiNiuTranslateMode:
		t = new(XunFeiNiuConfigType)
	case TencentTranslateMode:
		t = new(TencentConfigType)
	case HuoShanTranslateMode:
		t = new(HuoShanConfigType)
	case PaPaGoTranslateMode:
		t = new(PaPaGoConfigType)
	case FreeGoogleMode:
		t = new(FreeGoogle)
	default:
		err = errors.New("不支持的翻译")
		return
	}
	if config != nil {
		if err = gconv.Struct(config, t); err != nil {
			return
		}
	}
	return
}

// ChatGptTranslateMode ChatGPT 支持
const ChatGptTranslateMode = "ChatGPT"

// XunFeiTranslateMode 讯飞常用版本
const XunFeiTranslateMode = "XunFei"

// XunFeiNiuTranslateMode 讯飞新版
const XunFeiNiuTranslateMode = "XunFeiNiu"

const TencentTranslateMode = "Tencent"

// HuoShanTranslateMode 火山翻译
const HuoShanTranslateMode = "HuoShan"

// PaPaGoTranslateMode 啪啪GO翻译
const PaPaGoTranslateMode = "PaPaGo"

const (
	YouDaoTranslateMode = "YouDao"     // 有道
	BaiduTranslateMode  = "Baidu"      // 百度
	GoogleTranslateMode = "Google"     // 谷歌
	DeeplTranslateMode  = "Deepl"      // Deepl
	FreeGoogleMode      = "FreeGoogle" // 谷歌免费翻译
)

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
