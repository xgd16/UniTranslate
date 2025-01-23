package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
	"uniTranslate/src/buffer"
	"uniTranslate/src/command"
	"uniTranslate/src/devices"
	"uniTranslate/src/global"
	"uniTranslate/src/service"
	"uniTranslate/src/service/queue"
	"uniTranslate/src/translate"

	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/encoding/gjson"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/xgd16/gf-x-tool/xgraylog"
	"github.com/xgd16/gf-x-tool/xhttp"
	"github.com/xgd16/gf-x-tool/xlib"
	"github.com/xgd16/gf-x-tool/xmonitor"
)

func main() {
	// 创建命令
	mainCmd := &gcmd.Command{
		Name:        "main",
		Brief:       "开启 HTTP 服务",
		Description: "开启 HTTP API 服务",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			g.Log().Infof(ctx, "\n\n\n%s\nGitHub: %s\n文档地址: %s\n", gbase64.MustDecodeToString(global.LogoStr), "https://github.com/xgd16/UniTranslate", "https://uni-translate-doc.todream.net/")
			global.RunMode = global.HttpMode
			xmonitor.InitPrometheusMetric("uniTranslate", "uniTranslate")
			initHandler()
			// 初始化系统服务
			service.InitService()
			// 维持
			xlib.Maintain(nil)
			return
		},
	}
	translateCmd := &gcmd.Command{
		Name:        "translate",
		Brief:       "命令行翻译",
		Description: "开启命令行翻译",
		Arguments: []gcmd.Argument{
			{Name: "from", Brief: "源语言", IsArg: true},
			{Name: "to", Brief: "目标语言", IsArg: true},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			global.RunMode = global.CmdMode
			initHandler()
			go queue.Service()
			time.Sleep(700 * time.Millisecond)
			err = command.Translate(ctx, parser)
			return
		},
	}
	if err := mainCmd.AddCommand(translateCmd); err != nil {
		panic(err)
	}
	mainCmd.Run(gctx.New())
}

func initHandler() {
	runtime.SetMutexProfileFraction(1) // (非必需)开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // (非必需)开启对阻塞操作的跟踪
	// 初始化系统配置
	if err := global.InitSystemConfig(); err != nil {
		panic(err)
	}
	// 初始化基础
	baseInit()
}

func baseInit() {
	xhttp.RespErrorMsg = true
	// 初始化缓存
	global.GfCache = initCache()
	// 初始化记录驱动
	initRecordDevice()
	// 初始化配置驱动
	initConfigDevice()
	// 开启翻译支持
	translate.InitTranslate()
	// 初始化 chatGPT 需要的数据
	translate.ChatGPTLangConfig = gjson.MustEncodeString(translate.BaseTranslateConf[translate.ChatGptTranslateMode])
	// 初始化缓冲区
	if err := buffer.Buffer.Init(false); err != nil {
		panic(err)
	}
	// 配置 GrayLog 基础配置 host 和 port
	if global.SystemConfig.Get("graylog.open").Bool() {
		xgraylog.SetGrayLogConfig(global.SystemConfig.Get("graylog"))
		// 配置默认日志
		glog.SetDefaultHandler(xgraylog.SwitchToGraylog(global.ServerConfig.Name, func(ctx context.Context, m g.Map) {

		}))
	}
}

func initRecordDevice() {
	switch global.ServerConfig.RecordDeviceMode {
	case "sqlite":
		devices.RecordHandler = devices.NewSqlLiteRecordDevice()
	default:
		devices.RecordHandler = devices.NewMySQLRecordDevice()
	}
	if err := devices.RecordHandler.Init(global.GfCache, global.ServerConfig.CacheMode, global.ServerConfig.CachePlatform, global.ServerConfig.CacheRefreshOnStartup); err != nil {
		panic(fmt.Errorf("记录配置驱动初始化出错 %s", err))
	}
}

func initConfigDevice() {
	switch global.ServerConfig.ConfigDeviceMode {
	case "xdb":
		devices.ConfigDevice = devices.NewXDbConfigDevice()
	case "sqlite":
		devices.ConfigDevice = devices.NewSqlLiteConfigDevice()
	case "mysql":
		devices.ConfigDevice = devices.NewMySQLConfigDevice()
	default:
		devices.ConfigDevice = devices.NewXDbConfigDevice()
	}
	if err := devices.ConfigDevice.Init(); err != nil {
		panic(fmt.Errorf("翻译配置驱动初始化出错 %s", err))
	}
}

// 初始化 缓存
func initCache() *gcache.Cache {
	c := gcache.New()
	switch global.ServerConfig.CacheMode {
	case "redis":
		c.SetAdapter(gcache.NewAdapterRedis(g.Redis()))
	case "mem":
		c.SetAdapter(gcache.NewAdapterMemory())
	}
	return c
}
