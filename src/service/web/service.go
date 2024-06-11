package web

import (
	"uniTranslate/src/service/web/controller"
	"uniTranslate/src/service/web/route"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/xgd16/gf-x-tool/xmonitor"
)

func Service() {
	server := g.Server()
	// 路由注册
	server.Group("/api", route.Api)
	server.BindHandler("/ip", controller.GetIp)
	server.BindHandler("/metrics", xmonitor.PrometheusHttp)
	// 如果存在操作端的话启动静态地址
	if gfile.IsDir("./dist") {
		server.SetServerRoot("./dist")
	}
	// 基本配置
	server.BindMiddlewareDefault(ghttp.MiddlewareCORS)
	server.BindStatusHandler(404, func(r *ghttp.Request) {
		if gstr.Count(r.GetHeader("Content-Type"), "application/json") >= 1 {
			r.Response.WriteJsonExit(g.Map{
				"code": 1001,
				"msg":  "404 not find",
				"time": gtime.Now().UnixMilli(),
			})
		}
		r.Response.Writefln(`
			<div style="text-align:center;"><div style="font-size: 5rem">404</div><div style="font-size: 3rem">%s</div></div>
		`, gtime.Now().Format("Y-m-d H:i:s"))
	})
	server.EnablePProf()
	// 启动web服务
	server.Run()
}
