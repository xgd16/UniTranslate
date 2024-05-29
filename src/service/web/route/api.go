package route

import (
	"uniTranslate/src/service/web/controller"
	"uniTranslate/src/service/web/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Api 路由注册
func Api(r *ghttp.RouterGroup) {
	rA := r.Clone().Middleware(middleware.AuthVerifyMiddleware)
	rA.POST("/saveConfig", controller.SaveConfig)
	rA.GET("/getConfigList", controller.GetConfigList)
	rA.POST("/translate", controller.Translate)
	rA.POST("/aggregateTranslate", controller.AggregateTranslate)
	rA.GET("/refreshConfigCache", controller.RefreshConfigCache)
	rA.GET("/getCountRecord", controller.GetCountRecord)
	rA.GET("/getRequestRecord", controller.GetRequestRecord)
	r.GET("/getLangList", controller.GetLangList)
	r.GET("/getSystemInitConfig", controller.GetSystemInitConfig)
	rA.GET("/cleanCache", controller.CleanCache)
	rA.GET("/cacheSize", controller.CacheSize)
	rA.POST("/delConfig", controller.DelConfig)
	rA.POST("/updateStatus", controller.UpdateStatus)
}
