package route

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"uniTranslate/src/service/web/controller"
	"uniTranslate/src/service/web/middleware"
)

// Api 路由注册
func Api(r *ghttp.RouterGroup) {
	rA := r.Clone().Middleware(middleware.AuthVerifyMiddleware)
	// 请注意注册的微服务接口必须统一标准为 POST 方式
	rA.POST("/addConfig", controller.AddConfig)
	rA.POST("/translate", controller.Translate)
	r.GET("/getCountRecord", controller.GetCountRecord)
	r.GET("/getRequestRecord", controller.GetRequestRecord)
}
