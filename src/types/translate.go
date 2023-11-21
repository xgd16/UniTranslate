package types

import "github.com/gogf/gf/v2/os/gcache"

type TranslateData struct {
	Md5              string   `json:"-"`
	TranslateTextArr []string `json:"translate" orm:"translate"`
	From             string   `json:"from" orm:"fromLang"`
	To               string   `json:"to" orm:"toLang"`
	Platform         string   `json:"platform" orm:"platform"`
	OriginalText     string   `json:"originalText" orm:"text"`
	OriginalTextMd5  string   `json:"-" orm:"textMd5"`
	OriginalTextLen  int      `json:"originalTextLen" orm:"textLen"`
	TranslationLen   int      `json:"translationLen" orm:"translationLen"`
}

type StatisticsInterface interface {
	// Init 初始化数据库
	Init(cache *gcache.Cache, cacheMode string, cachePlatform bool) error
	// CountRecord 计数统计
	CountRecord(data *CountRecordData) error
	// RequestRecord 请求记录
	RequestRecord(data *RequestRecordData) error
	// CreateEvent 触发创建事件
	CreateEvent(data *TranslatePlatform) error
	// SaveCache 存储翻译结果到缓存
	SaveCache(data *SaveData) error
	// GetterCache 获取翻译结果
	GetterCache(fn func(data []*TranslateData) (err error)) error
}

type CountRecordData struct {
	Data *TranslateData `json:"data"`
	Ok   bool           `json:"ok"`
}

type RequestRecordData struct {
	ClientIp string `json:"clientIp"`
	Body     string `json:"body"`
	Time     int64  `json:"time"`
	Ok       bool   `json:"ok"`
	ErrMsg   error  `json:"errMsg"`
}

type SaveData struct {
	Data *TranslateData `json:"data"`
}
