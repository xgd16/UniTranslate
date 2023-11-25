package main

import (
	"fmt"
	"uniTranslate/src/global"
	"uniTranslate/src/service"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/xgd16/gf-x-tool/xgraylog"
	"github.com/xgd16/gf-x-tool/xhttp"
	"github.com/xgd16/gf-x-tool/xlib"
	"github.com/xgd16/gf-x-tool/xtranslate"
)

func main() {
	// 初始化系统配置
	global.InitSystemConfig()
	// 初始化基础
	baseInit()
	// 初始化系统服务
	service.InitService()
	// 维持
	xlib.Maintain(nil)
}

func baseInit() {
	xhttp.RespErrorMsg = true
	// 开启翻译支持
	xtranslate.InitTranslate()
	// 初始化缓冲区
	if err := global.Buffer.Init(); err != nil {
		panic(err)
	}
	// 初始化缓存
	global.GfCache = initCache()
	// 初始化数据库
	if err := initDataBase(); err != nil {
		panic(fmt.Errorf("初始化数据库失败 %s", err))
	}
	// 配置 GrayLog 基础配置 host 和 port
	if global.SystemConfig.Get("grayLog.open").Bool() {
		xgraylog.SetGrayLogConfig(
			global.SystemConfig.Get("grayLog.host").String(),
			global.SystemConfig.Get("grayLog.port").Int(),
		)
		// 配置默认日志
		name := global.SystemConfig.Get("server.name").String()
		glog.SetDefaultHandler(xgraylog.SwitchToGraylog(name))
	}
}

// 初始化 数据库信息
func initDataBase() (err error) {
	err = global.StatisticalProcess.Init(global.GfCache, global.CacheMode, global.CachePlatform, global.CacheRefreshOnStartup)
	return
}

// 初始化 缓存
func initCache() *gcache.Cache {
	c := gcache.New()
	switch global.CacheMode {
	case "redis":
		c.SetAdapter(gcache.NewAdapterRedis(g.Redis()))
		break
	case "mem":
		c.SetAdapter(gcache.NewAdapterMemory())
		break
	}
	return c
}
