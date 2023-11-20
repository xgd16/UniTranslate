package types

type TranslateData struct {
	Md5              string   `json:"-"`
	TranslateTextArr []string `json:"translate"`
	From             string   `json:"from"`
	To               string   `json:"to"`
	Platform         string   `json:"platform"`
	OriginalText     string   `json:"originalText"`
	OriginalTextLen  int      `json:"originalTextLen"`
	TranslationLen   int      `json:"translationLen"`
}

type StatisticsInterface interface {
	// CountRecord 计数统计
	CountRecord(data *CountRecordData) error
	// RequestRecord 请求记录
	RequestRecord(data *RequestRecordData) error
	// CreateEvent 触发创建事件
	CreateEvent(data *TranslatePlatform) error
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
