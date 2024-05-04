package types

import (
	"uniTranslate/src/translate"

	"github.com/gogf/gf/v2/os/gcache"
)

// TranslateData 翻译数据
type TranslateData struct {
	Md5             string                     `json:"-"`
	Translate       []*translate.TranslateResp `json:"translate"`
	To              string                     `json:"to" orm:"toLang"`
	Platform        string                     `json:"platform" orm:"platform"`
	OriginalText    []string                   `json:"originalText" orm:"text"`
	OriginalTextStr *string                    `json:"-" orm:"-"`
	OriginalTextMd5 string                     `json:"-" orm:"textMd5"`
	OriginalTextLen int                        `json:"originalTextLen" orm:"textLen"`
}

// CountRecord 计数记录
type CountRecord struct {
	SerialNumber string `json:"serialNumber"`
	SuccessCount int    `json:"successCount"`
	ErrorCount   int    `json:"errorCount"`
	CharCount    int    `json:"charCount"`
}

// StatisticsInterface 统计接口
type StatisticsInterface interface {
	// Init 初始化数据库
	Init(cache *gcache.Cache, cacheMode string, cachePlatform, cacheRefreshOnStartup bool) error
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
	// GetCountRecord 获取计数记录
	GetCountRecord() (data map[string]*CountRecord, err error)
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
	Platform string `json:"platform"`
	ErrMsg   error  `json:"errMsg"`
	TraceId  string `json:"traceId" orm:"tId"`
	TakeTime int    `json:"takeTime"`
}

type SaveData struct {
	Data *TranslateData `json:"data"`
}
