package controller

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xgd16/gf-x-tool/x"
	"uniTranslate/src/global"
)

func GetCountRecord(r *ghttp.Request) {
	data, err := g.Model("count_record").OrderDesc("id").All()
	for i, item := range data {
		data[i]["name"] = global.XDB.GetGJson().Get(fmt.Sprintf("translate.%s.platform", item["serialNumber"].String()))
	}
	x.FastResp(r, err).Resp()
	x.FastResp(r).SetData(data).Resp()
}

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
