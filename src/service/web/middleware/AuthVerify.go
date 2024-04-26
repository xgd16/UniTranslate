package middleware

import (
	"uniTranslate/src/global"
	"uniTranslate/src/lib"

	"github.com/gogf/gf/v2/net/ghttp"
)

// AuthVerifyMiddleware 身份验证
func AuthVerifyMiddleware(r *ghttp.Request) {
	authKey := r.Header.Get("auth_key")
	// 根据设置的验证方式验证请求
	var pass bool
	if global.KeyMode == 1 {
		if authKey == "" {
			authKey = r.Get("key").String()
		}
		pass = authKey != global.ServiceKey
	} else {
		pass = authKey != lib.AuthEncrypt(global.ServiceKey, r.GetMap())
	}

	if pass {
		r.Response.Status = 404
		r.Response.WriteExit()
	}

	r.Middleware.Next()
}
