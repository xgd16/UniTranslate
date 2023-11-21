package types

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"math"
)

type MySqlStatistics struct{}

func (m *MySqlStatistics) Init(cache *gcache.Cache, cacheMode string, cachePlatform bool) (err error) {
	ctx := gctx.New()
	// 初始化基础结构数据
	type MySqlInit struct {
		TableName string
		Table     string
		Index     []string
	}
	// 创建表信息
	initData := []*MySqlInit{
		{
			TableName: "count_record",
			Table:     "CREATE TABLE count_record ( id int UNSIGNED PRIMARY KEY AUTO_INCREMENT, serialNumber varchar(255) NOT NULL, successCount int NULL DEFAULT 0, errorCount int NULL DEFAULT 0, charCount int NULL DEFAULT 0, createTime datetime(6) NOT NULL, updateTime datetime(6) NULL );",
			Index: []string{
				"CREATE INDEX count_record_charCount_index ON count_record (charCount);",
				"CREATE INDEX count_record_createTime_index ON count_record (createTime);",
				"CREATE INDEX count_record_errorCount_index ON count_record (errorCount);",
				"CREATE INDEX count_record_serialNumber_index ON count_record (serialNumber);",
				"CREATE INDEX count_record_successCount_index ON count_record (successCount);",
			},
		},
		{
			TableName: "request_record",
			Table:     "CREATE TABLE request_record ( id int UNSIGNED PRIMARY KEY AUTO_INCREMENT, clientIp varchar(255) NULL, body text NULL, status tinyint(1) NULL, errMsg text NULL, createTime datetime(6) NULL, updateTime datetime(6) NULL );",
			Index: []string{
				"CREATE INDEX request_record_clientIp_index ON request_record (clientIp);",
				"CREATE INDEX request_record_createTime_index ON request_record (createTime);",
				"CREATE INDEX request_record_status_index ON request_record (status);",
			},
		},
		{
			TableName: "translate_cache",
			Table:     "CREATE TABLE translate_cache ( id bigint UNSIGNED AUTO_INCREMENT, translate json NOT NULL COMMENT '翻译后结果', text text NOT NULL, fromLang varchar(16) NULL COMMENT '翻译前语言', toLang varchar(16) NULL COMMENT '翻译后语言', textMd5 char(32) NOT NULL COMMENT '翻译前语言md5值', platform varchar(16) NOT NULL COMMENT '翻译平台 Baidu YouDao Google Deepl', textLen int NULL DEFAULT 0 COMMENT '原文文字长度', translationLen int NULL DEFAULT 0 COMMENT '翻译后文字长度', createTime datetime NOT NULL, updateTime datetime NULL, PRIMARY KEY (id) );",
			Index: []string{
				"CREATE INDEX translate_cache_createTime_index ON translate_cache (createTime);",
				"CREATE INDEX translate_cache_fromLang_index ON translate_cache (fromLang);",
				"CREATE INDEX translate_cache_platform_index ON translate_cache (platform);",
				"CREATE INDEX translate_cache_textLen_index ON translate_cache (textLen);",
				"CREATE INDEX translate_cache_toLang_index ON translate_cache (toLang);",
				"CREATE INDEX translate_cache_translationLen_index ON translate_cache (translationLen);",
				"CREATE INDEX translate_cache_textMd5_index ON translate_cache (textMd5);",
			},
		},
	}
	// 循环处理
	var hasTable *gvar.Var
	for _, item := range initData {
		// 判断是否存在
		hasTable, err = g.DB().GetValue(ctx, fmt.Sprintf("SHOW TABLES LIKE '%s'", item.TableName))
		if err != nil {
			return
		}
		if hasTable.String() == item.TableName {
			continue
		}
		// 执行创建
		if _, err = g.DB().Exec(ctx, item.Table); err != nil {
			return
		}
		for _, index := range item.Index {
			if _, err = g.DB().Exec(ctx, index); err != nil {
				return
			}
		}
	}
	// 存储到缓存
	err = saveToCache(ctx, cache, m, cacheMode, cachePlatform)
	return
}

// 存储到缓存
func saveToCache(ctx context.Context, cache *gcache.Cache, m *MySqlStatistics, cacheMode string, cachePlatform bool) (err error) {
	const keyName = "Translate-"
	if cacheMode == "redis" {
		conn, connErr := g.Redis().Conn(ctx)
		if connErr != nil {
			return connErr
		}
		defer func() {
			if err = conn.Close(ctx); err != nil {
				return
			}
		}()
		keys, keysErr := conn.Do(ctx, "KEYS", fmt.Sprintf("%s*", keyName))
		if keysErr != nil {
			return keysErr
		}
		for _, item := range keys.Strings() {
			if _, delErr := conn.Do(ctx, "DEL", item); delErr != nil {
				return delErr
			}
		}
	}
	// 内存缓存是否 包含 平台
	err = m.GetterCache(func(data []*TranslateData) (err error) {
		for _, item := range data {
			var md5 string
			if cachePlatform {
				md5 = gmd5.MustEncrypt(fmt.Sprintf("to:%s-text:%s-platform:%s", item.To, item.OriginalText, item.Platform))
			} else {
				md5 = gmd5.MustEncrypt(fmt.Sprintf("to:%s-text:%s", item.To, item.OriginalText))
			}
			if err = cache.Set(ctx, fmt.Sprintf("%s%s", keyName, md5), item, 0); err != nil {
				return
			}
		}
		return
	})
	return
}

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

func (m *MySqlStatistics) SaveCache(data *SaveData) error {
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

func (m *MySqlStatistics) GetterCache(fn func(data []*TranslateData) (err error)) error {
	pageSize := 1
	model := g.Model("translate_cache")
	dbSize, err := model.Clone().Count()
	if err != nil {
		return err
	}
	for i := 0; i < int(math.Ceil(float64(dbSize)/float64(pageSize))); i++ {
		newData := make([]*TranslateData, 0)
		if err = model.Clone().Page(i+1, pageSize).Scan(&newData); err != nil {
			return err
		}
		if err = fn(newData); err != nil {
			return err
		}
	}
	return nil
}
