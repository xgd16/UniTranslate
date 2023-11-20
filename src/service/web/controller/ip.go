package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xgd16/gf-x-tool/x"
)

func GetIp(r *ghttp.Request) {
	x.FastResp(r).SetData(g.Map{
		"ip": r.GetClientIp(),
	}).Resp()
}
