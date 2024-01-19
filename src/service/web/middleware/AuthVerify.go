package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"uniTranslate/src/global"
	"uniTranslate/src/lib"
)

// AuthVerifyMiddleware 身份验证
func AuthVerifyMiddleware(r *ghttp.Request) {
	authKey := r.Header.Get("auth_key")

	if authKey != lib.AuthEncrypt(global.ServiceKey, r.GetMap()) {
		r.Response.Status = 404
		r.Response.WriteExit()
	}

	r.Middleware.Next()
}
