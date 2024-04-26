package route

import (
	"uniTranslate/src/service/web/controller"
	"uniTranslate/src/service/web/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Api 路由注册
func Api(r *ghttp.RouterGroup) {
	rA := r.Clone().Middleware(middleware.AuthVerifyMiddleware)
	rA.POST("/addConfig", controller.AddConfig)
	rA.GET("/getConfigList", controller.GetConfigList)
	rA.POST("/translate", controller.Translate)
	rA.GET("/refreshConfigCache", controller.RefreshConfigCache)
	r.GET("/getCountRecord", controller.GetCountRecord)
	r.GET("/getRequestRecord", controller.GetRequestRecord)
	r.GET("/getLangList", controller.GetLangList)
	r.GET("/getSystemInitConfig", controller.GetSystemInitConfig)
	r.POST("/delConfig", controller.DelConfig)
	r.POST("/updateStatus", controller.UpdateStatus)
}
