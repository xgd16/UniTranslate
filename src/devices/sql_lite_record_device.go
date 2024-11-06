package devices

import (
	"errors"
	"math"
	"strings"
	"uniTranslate/src/global"
	"uniTranslate/src/lib"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
)

type SqlLiteRecordDevice struct{}

func NewSqlLiteRecordDevice() *SqlLiteRecordDevice {
	return &SqlLiteRecordDevice{}
}

func (t *SqlLiteRecordDevice) Init(cache *gcache.Cache, cacheMode string, cachePlatform, cacheRefreshOnStartup bool) (err error) {
	initData := []*types.SQLInitItem{
		{
			TableName: "request_record",
			Table:     "create table request_record( id integer primary key autoincrement, tId varchar(255) null, clientIp varchar(255) null, body text null, status tinyint null, errMsg text null, takeTime int null, platform varchar(64) null, cacheId varchar(255), createTime datetime null, updateTime datetime null);",
			Index: []string{
				"create index request_record_clientIp_index on request_record (clientIp);",
				"create index request_record_createTime_index on request_record (createTime);",
				"create index request_record_platform_index on request_record (platform);",
				"create index request_record_status_index on request_record (status);",
				"create index request_record_tId_index on request_record (tId);",
				"create index request_record_takeTime_index on request_record (takeTime);",
				"create index request_record_cacheId_index on request_record (cacheId);",
			},
		},
		{
			TableName: "count_record",
			Table:     "create table count_record( id integer primary key autoincrement, serialNumber varchar(255) not null, successCount int default 0 null, errorCount int default 0 null, charCount bigint default 0 null, createTime datetime not null, updateTime datetime null);",
			Index: []string{
				"create index count_record_charCount_index on count_record (charCount);",
				"create index count_record_createTime_index on count_record (createTime);",
				"create index count_record_errorCount_index on count_record (errorCount);",
				"create index count_record_serialNumber_index on count_record (serialNumber);",
				"create index count_record_successCount_index on count_record (successCount);",
			},
		},
		{
			TableName: "translate_cache",
			Table:     "create table translate_cache( id integer primary key autoincrement, translate text not null, text text not null, textMd5 char(32) not null, fromLang varchar(16) null, toLang varchar(16) null, platform varchar(16) not null, textLen int default 0 null, translationLen int default 0 null, cacheId varchar(255), createTime datetime not null, updateTime datetime null);",
			Index: []string{
				"create index translate_cache_createTime_index on translate_cache (createTime);",
				"create index translate_cache_fromLang_index on translate_cache (fromLang);",
				"create index translate_cache_platform_index on translate_cache (platform);",
				"create index translate_cache_textLen_index on translate_cache (textLen);",
				"create index translate_cache_textMd5_index on translate_cache (textMd5);",
				"create index translate_cache_toLang_index on translate_cache (toLang);",
				"create index translate_cache_translationLen_index on translate_cache (translationLen);",
				"create index translate_cache_cacheId_index on translate_cache (cacheId);",
			},
		},
	}
	db := g.DB()
	var exists bool
	for _, item := range initData {
		exists, err = lib.SqliteTableIsExists(db, item.TableName)
		if err != nil {
			return
		}
		if exists {
			continue
		}
		if _, err = db.Exec(global.Ctx, item.Table); err != nil {
			return
		}
		for _, item1 := range item.Index {
			if _, err = db.Exec(global.Ctx, item1); err != nil {
				return
			}
		}
	}
	return
}

func (t *SqlLiteRecordDevice) CountRecord(data *types.CountRecordData) error {
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
	if data.Ok {
		_, err = model.Clone().Increment("charCount", data.Data.OriginalTextLen)
	}
	return err
}

func (t *SqlLiteRecordDevice) RequestRecord(data *types.RequestRecordData) error {
	var errMsg string
	if data.ErrMsg != nil {
		errMsg = data.ErrMsg.Error()
	}
	_, err := g.Model("request_record").Data(g.Map{
		"clientIp": data.ClientIp,
		"body":     data.Body,
		"status":   gconv.Int(data.Ok),
		"errMsg":   errMsg,
		"takeTime": data.TakeTime,
		"platform": data.Platform,
		"tId":      data.TraceId,
		"cacheId":  data.CacheId,
	}).Insert()
	return err
}

func (t *SqlLiteRecordDevice) CreateEvent(data *types.TranslatePlatform) error {
	_, err := g.Model("count_record").Data(g.Map{
		"serialNumber": data.Md5,
	}).Insert()
	return err
}

func (t *SqlLiteRecordDevice) SaveCache(data *types.SaveData) error {
	if data == nil || data.Data == nil {
		return nil
	}
	count, err := g.Model("translate_cache").Count(g.Map{
		"textMd5":  data.Data.OriginalTextMd5,
		"platform": data.Data.Platform,
	})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	if _, err = g.Model("translate_cache").Insert(data.Data); err != nil {
		return err
	}
	return nil
}

func (t *SqlLiteRecordDevice) GetterCache(fn func(data []*types.TranslateData) (err error)) error {
	pageSize := 1
	model := g.Model("translate_cache")
	dbSize, err := model.Clone().Count()
	if err != nil {
		return err
	}
	for i := 0; i < int(math.Ceil(float64(dbSize)/float64(pageSize))); i++ {
		newData := make([]*types.TranslateData, 0)
		if err = model.Clone().Page(i+1, pageSize).Scan(&newData); err != nil {
			return err
		}
		for _, item := range newData {
			textStr := strings.Join(item.OriginalText, "\n")
			item.OriginalTextStr = &textStr
		}
		if err = fn(newData); err != nil {
			return err
		}
	}
	return nil
}

func (t *SqlLiteRecordDevice) GetCountRecord() (data map[string]*types.CountRecord, err error) {
	all, err := g.Model("count_record").All()
	if err != nil {
		return
	}
	if all.IsEmpty() {
		return
	}
	temp := make([]*types.CountRecord, 0)
	if err = all.Structs(&temp); err != nil {
		return
	}
	data = make(map[string]*types.CountRecord)
	for _, item := range temp {
		data[item.SerialNumber] = item
	}
	return
}
