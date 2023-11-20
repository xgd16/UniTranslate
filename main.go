package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/xgd16/gf-x-tool/xgraylog"
	"github.com/xgd16/gf-x-tool/xhttp"
	"github.com/xgd16/gf-x-tool/xtranslate"
	"uniTranslate/src/global"
	"uniTranslate/src/lib"
	"uniTranslate/src/service"
)

func main() {
	// 初始化系统配置
	global.InitSystemConfig()
	// 初始化基础
	baseInit()
	// 初始化系统服务
	service.InitService()
	// 维持
	lib.Maintain()
}

func baseInit() {
	xhttp.RespErrorMsg = true
	// 开启翻译支持
	xtranslate.InitTranslate()
	// 初始化缓冲区
	if err := global.Buffer.Init(); err != nil {
		panic(err)
	}
	// 配置 GrayLog 基础配置 host 和 port
	if !global.SystemConfig.Get("server.devMode").Bool() {
		xgraylog.SetGrayLogConfig(
			global.SystemConfig.Get("grayLog.host").String(),
			global.SystemConfig.Get("grayLog.port").Int(),
		)
		// 配置默认日志
		name := global.SystemConfig.Get("server.name").String()
		glog.SetDefaultHandler(xgraylog.SwitchToGraylog(name))
	}
}
