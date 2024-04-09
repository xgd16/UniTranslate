package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xgd16/gf-x-tool/x"
	"uniTranslate/src/global"
)

type SystemInitConfigData struct {
	AuthMode int `json:"authMode"` // key mode
}

// GetSystemInitConfig 获取系统初始化配置
func GetSystemInitConfig(r *ghttp.Request) {

	x.FastResp(r).Resp("", &SystemInitConfigData{
		AuthMode: global.KeyMode,
	})
}
