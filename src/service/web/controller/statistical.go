package controller

import (
	"uniTranslate/src/global"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xgd16/gf-x-tool/x"
)

// GetCountRecord 获取账号计数信息
func GetCountRecord(r *ghttp.Request) {
	data, err := g.Model("count_record").OrderDesc("id").All()
	x.FastResp(r, err).Resp()
	device, err := global.GetConfigDevice()
	x.FastResp(r, err).Resp()
	for i, item := range data {
		var (
			info *types.TranslatePlatform
			ok   bool
		)
		info, ok, err = device.GetTranslateInfo(item["serialNumber"].String())
		x.FastResp(r, err).Resp()
		if !ok {
			continue
		}
		data[i]["name"] = gvar.New(info.Platform)
	}
	x.FastResp(r, err).Resp()
	x.FastResp(r).SetData(data).Resp()
}

// GetRequestRecord 获取账号请求信息
func GetRequestRecord(r *ghttp.Request) {
	page := r.Get("page", 1).Int()
	size := r.Get("size", 10).Int()

	m := g.Model("request_record")

	data, err := m.Clone().Page(page, size).OrderDesc("id").All()
	x.FastResp(r, err).Resp()
	count, err := m.Clone().Count()
	x.FastResp(r, err).Resp()
	x.FastResp(r).SetData(g.Map{
		"list":  data,
		"count": count,
	}).Resp()
}
