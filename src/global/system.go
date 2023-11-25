package global

import (
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/xgd16/gf-x-tool/xstorage"
)

// SystemConfig 系统配置信息
var SystemConfig *gjson.Json

var CacheMode = "mem"
var CachePlatform = false

// InitSystemConfig 初始化系统配置信息
func InitSystemConfig() {
	cfg, err := g.Cfg().Data(gctx.New())
	if err != nil {
		panic("初始化系统配置错误: " + err.Error())
	}
	SystemConfig = gjson.New(cfg, true)
	// 初始化配置的缓存模式
	CacheMode = SystemConfig.Get("server.cacheMode").String()
	CachePlatform = SystemConfig.Get("server.cachePlatform").Bool()
	CacheRefreshOnStartup = SystemConfig.Get("server.cacheRefreshOnStartup").Bool()
}

// XDB 文件式存储
var XDB = xstorage.CreateXDB()

// CacheRefreshOnStartup 启动时是否从数据库刷新缓存 (会先清除缓存里所有的 缓存 在从数据库逐条初始化 数据 慎用!!!)
var CacheRefreshOnStartup = false

// GfCache 全局缓存
var GfCache *gcache.Cache

var StatisticalProcess types.StatisticsInterface = new(types.MySqlStatistics)
