package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"uniTranslate/src/global"
)

func AuthVerifyMiddleware(r *ghttp.Request) {
	if r.Get("key").String() != global.SystemConfig.Get("server.key").String() {
		r.Response.Status = 404
		r.Response.WriteExit()
	}

	r.Middleware.Next()
}
