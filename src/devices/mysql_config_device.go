package devices

import (
	"context"
	"fmt"
	"uniTranslate/src/global"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/frame/gins"
)

type MySQLConfigDevice struct {
	dbName    string
	tableName string
}

func NewMySQLConfigDevice() *MySQLConfigDevice {
	return &MySQLConfigDevice{
		global.ConfigDeviceMySqlDb,
		"translation_platform",
	}
}

func (t *MySQLConfigDevice) Init() (err error) {
	ctx := context.Background()
	// 创建表信息
	initData := []*types.MySQLInitItem{
		{
			TableName: "translation_platform",
			Table:     "CREATE TABLE translation_platform ( id int UNSIGNED PRIMARY KEY AUTO_INCREMENT, md5 char(32) NOT NULL, translatedPlatform varchar(255) NOT NULL COMMENT '翻译平台', useBytes int NOT NULL DEFAULT 0 COMMENT '使用字节', errorTimes int NULL DEFAULT 0 COMMENT '报错次数', status tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态0:关闭1:开启', translationLevel int NULL COMMENT '翻译等级', createTime datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间', updateTime datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间', cfg json NOT NULL COMMENT '配置', type varchar(15) NULL COMMENT '翻译平台类型 YouDao Baidu' ) CHARSET = utf8mb4 COMMENT '翻译平台';",
			Index: []string{
				"CREATE INDEX createTime_index ON translation_platform (createTime);",
				"CREATE INDEX status_index ON translation_platform (status);",
				"CREATE INDEX translatedPlatform_index ON translation_platform (translatedPlatform);",
				"CREATE INDEX translation_platform_md5_index ON translation_platform (md5);",
			},
		},
	}
	// 循环处理
	var hasTable *gvar.Var
	for _, item := range initData {
		// 判断是否存在
		hasTable, err = t.db().GetValue(ctx, fmt.Sprintf("SHOW TABLES LIKE '%s'", item.TableName))
		if err != nil {
			return
		}
		if hasTable.String() == item.TableName {
			if err = t.initMd5(); err != nil {
				return
			}
			continue
		}
		// 执行创建
		if _, err = t.db().Exec(ctx, item.Table); err != nil {
			return
		}
		for _, index := range item.Index {
			if _, err = t.db().Exec(ctx, index); err != nil {
				return
			}
		}
	}
	// 初始化计数表
	if err = t.initCountRecord(); err != nil {
		return
	}
	return
}

func (t *MySQLConfigDevice) GetConfig(refresh bool) (mapData map[string]*types.TranslatePlatform, err error) {
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

func (t *MySQLConfigDevice) GetTranslateInfo(serialNumber string) (platform *types.TranslatePlatform, ok bool, err error) {
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

func (t *MySQLConfigDevice) SaveConfig(serialNumber string, isUpdate bool, data *types.TranslatePlatform) (err error) {
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

func (t *MySQLConfigDevice) DelConfig(serialNumber string) (err error) {
	if _, err = t.model().Where("md5", serialNumber).Delete(); err != nil {
		return
	}
	return
}

func (t *MySQLConfigDevice) UpdateStatus(serialNumber string, status int) (err error) {
	if _, err = t.model().Update(g.Map{
		"status": status,
	}, g.Map{
		"md5": serialNumber,
	}); err != nil {
		return
	}
	return
}

func (t *MySQLConfigDevice) db() gdb.DB {
	return gins.Database(t.dbName)
}

func (t *MySQLConfigDevice) model() *gdb.Model {
	return t.db().Model(t.tableName)
}

func (t *MySQLConfigDevice) initMd5() (err error) {
	allT, err := t.model().All()
	if err != nil {
		return
	}
	all := make([]*types.TranslatePlatform, 0)
	if err = allT.Structs(&all); err != nil {
		return
	}
	for _, item := range all {
		if item.Md5 != "" {
			continue
		}
		item.InitMd5()
		if _, err = t.model().Update(g.Map{
			"md5": item.Md5,
		}, g.Map{
			"id": item.Id,
		}); err != nil {
			return
		}
	}
	return
}

func (t *MySQLConfigDevice) initCountRecord() (err error) {
	all, err := t.model().Fields("md5").All()
	if err != nil {
		return
	}

	m := g.Model("count_record")

	for _, v := range all {
		md5 := v["md5"].String()
		count, err1 := m.Clone().Where("serialNumber", md5).Count()
		if err1 != nil {
			err = err1
			return
		}
		if count > 0 {
			continue
		}
		if _, err = m.Clone().InsertIgnore(g.Map{
			"serialNumber": md5,
		}); err != nil {
			return
		}
	}
	return
}
