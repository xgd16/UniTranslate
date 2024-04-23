package translate

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
