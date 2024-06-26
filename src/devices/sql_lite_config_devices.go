package devices

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/frame/gins"
	"uniTranslate/src/global"
	"uniTranslate/src/lib"
	"uniTranslate/src/types"
)

type SqlLiteConfigDevice struct {
	dbName    string
	tableName string
}

func NewSqlLiteConfigDevice() *SqlLiteConfigDevice {
	return &SqlLiteConfigDevice{
		global.ConfigDeviceDb,
		"translation_platform",
	}
}

func (t *SqlLiteConfigDevice) Init() (err error) {
	initData := []*types.SQLInitItem{
		{
			TableName: "translation_platform",
			Table:     "create table translation_platform( id integer primary key autoincrement, md5 varchar(32) not null, translatedPlatform varchar(255) not null, useBytes int default 0 not null, errorTimes int default 0 null, status tinyint default 1 not null, translationLevel int null, createTime datetime default CURRENT_TIMESTAMP not null, updateTime datetime default CURRENT_TIMESTAMP not null, cfg text not null, type varchar(15) null);",
			Index: []string{
				"create index createTime_index on translation_platform(createTime);",
				"create index status_index on translation_platform (status);",
				"create index translatedPlatform_index on translation_platform (translatedPlatform);",
			},
		},
	}
	var exists bool
	for _, item := range initData {
		exists, err = lib.SqliteTableIsExists(t.db(), item.TableName)
		if err != nil {
			return
		}
		if exists {
			continue
		}
		if _, err = t.db().Exec(global.Ctx, item.Table); err != nil {
			return
		}
		for _, item1 := range item.Index {
			if _, err = t.db().Exec(global.Ctx, item1); err != nil {
				return
			}
		}
	}
	return
}

func (t *SqlLiteConfigDevice) GetConfig(refresh bool) (mapData map[string]*types.TranslatePlatform, err error) {
	mapData = make(map[string]*types.TranslatePlatform)
	allT, err := t.model().All()
	if err != nil {
		return
	}
	var all []*types.TranslatePlatform
	if err = allT.Structs(&all); err != nil {
		return
	}
	for _, item := range all {
		mapData[item.Md5] = item
	}
	return
}

func (t *SqlLiteConfigDevice) db() gdb.DB {
	return gins.Database(t.dbName)
}

func (t *SqlLiteConfigDevice) model() *gdb.Model {
	return t.db().Model(t.tableName)
}

func (t *SqlLiteConfigDevice) GetTranslateInfo(serialNumber string) (platform *types.TranslatePlatform, ok bool, err error) {
	one, err := t.model().Where("md5", serialNumber).One()
	if err != nil {
		return
	}
	if err = one.Struct(&platform); err != nil {
		return
	}
	ok = !one.IsEmpty()
	return
}

func (t *SqlLiteConfigDevice) SaveConfig(serialNumber string, isUpdate bool, data *types.TranslatePlatform) (err error) {
	dataMap := g.Map{
		"translatedPlatform": data.Platform,
		"status":             data.Status,
		"translationLevel":   data.Level,
		"cfg":                data.Cfg,
		"type":               data.Type,
		"md5":                data.Md5,
	}
	if isUpdate {
		_, err = t.model().Where("md5", serialNumber).Update(dataMap)
		return
	}
	_, err = t.model().Insert(dataMap)
	return
}

func (t *SqlLiteConfigDevice) DelConfig(serialNumber string) (err error) {
	if _, err = t.model().Where("md5", serialNumber).Delete(); err != nil {
		return
	}
	return
}

func (t *SqlLiteConfigDevice) UpdateStatus(serialNumber string, status int) (err error) {
	if _, err = t.model().Update(g.Map{
		"status": status,
	}, g.Map{
		"md5": serialNumber,
	}); err != nil {
		return
	}
	return
}
