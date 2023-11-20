package global

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/xgd16/gf-x-tool/xstorage"
	"uniTranslate/src/types"
)

// SystemConfig 系统配置信息
var SystemConfig *gjson.Json

var CacheMode string = "mem"

// InitSystemConfig 初始化系统配置信息
func InitSystemConfig() {
	cfg, err := g.Cfg().Data(gctx.New())
	if err != nil {
		panic("初始化系统配置错误: " + err.Error())
	}
	SystemConfig = gjson.New(cfg, true)
	// 初始化配置的缓存模式
	CacheMode = SystemConfig.Get("server.cacheMode").String()
}

// XDB 文件式存储
var XDB = xstorage.CreateXDB()

// GfCache 全局缓存
var GfCache *gcache.Cache

var StatisticalProcess types.StatisticsInterface = new(types.MySqlStatistics)
