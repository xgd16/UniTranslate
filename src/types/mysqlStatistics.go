package types

import (
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type MySqlStatistics struct{}

func (m *MySqlStatistics) CountRecord(data *CountRecordData) error {
	if data.Data == nil {
		return errors.New("翻译参数异常")
	}
	model := g.Model("count_record").Where("serialNumber", data.Data.Md5)
	_, err := model.Clone().Increment(func() string {
		if data.Ok {
			return "successCount"
		} else {
			return "errorCount"
		}
	}(), 1)
	if err != nil {
		return err
	}
	_, err = model.Clone().Increment("charCount", data.Data.OriginalTextLen)
	return err
}

func (m *MySqlStatistics) RequestRecord(data *RequestRecordData) error {
	_, err := g.Model("request_record").Data(g.Map{
		"clientIp": data.ClientIp,
		"body":     data.Body,
		"status":   gconv.Int(data.Ok),
		"errMsg":   data.ErrMsg,
	}).Insert()
	return err
}

func (m *MySqlStatistics) CreateEvent(data *TranslatePlatform) error {
	_, err := g.Model("count_record").Data(g.Map{
		"serialNumber": data.md5,
	}).Insert()
	return err
}
