package main

import (
	"context"
	"fmt"
	"uniTranslate/src/buffer"
	"uniTranslate/src/devices"
	"uniTranslate/src/global"
	"uniTranslate/src/service"
	"uniTranslate/src/translate"

	"github.com/gogf/gf/v2/encoding/gjson"

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
	// 初始化翻译配置获取
	initTranslateConfigDevice()
	// 开启翻译支持
	xtranslate.InitTranslate()
	// 初始化 chatGPT 需要的数据
	global.ChatGPTLangConfig = gjson.MustEncodeString(xtranslate.BaseTranslateConf[translate.ChatGptTranslateMode])
	// 初始化缓冲区
	if err := buffer.Buffer.Init(); err != nil {
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
		glog.SetDefaultHandler(xgraylog.SwitchToGraylog(name, func(ctx context.Context, m g.Map) {

		}))
	}
}

// 初始化 数据库信息
func initDataBase() (err error) {
	err = global.StatisticalProcess.Init(global.GfCache, global.CacheMode, global.CachePlatform, global.CacheRefreshOnStartup)
	return
}

func initTranslateConfigDevice() {
	switch global.ConfigDeviceMode {
	case "xdb":
		global.ConfigDevice = devices.NewXDbConfigDevice()
	case "mysql":
		global.ConfigDevice = devices.NewMySQLConfigDevice()
	default:
		global.ConfigDevice = devices.NewXDbConfigDevice()
	}
	if err := global.ConfigDevice.Init(); err != nil {
		panic(fmt.Errorf("翻译配置驱动初始化出错 %s", err))
	}
}

// 初始化 缓存
func initCache() *gcache.Cache {
	c := gcache.New()
	switch global.CacheMode {
	case "redis":
		c.SetAdapter(gcache.NewAdapterRedis(g.Redis()))
	case "mem":
		c.SetAdapter(gcache.NewAdapterMemory())
	}
	return c
}
