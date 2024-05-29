package controller

import (
	"uniTranslate/src/buffer"
	"uniTranslate/src/global"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xgd16/gf-x-tool/x"
)

type SystemInitConfigData struct {
	AuthMode   int  `json:"authMode"`
	EditConfig bool `json:"editConfig"`
}

// GetSystemInitConfig 获取系统初始化配置
func GetSystemInitConfig(r *ghttp.Request) {
	x.FastResp(r).Resp("", &SystemInitConfigData{
		AuthMode:   global.KeyMode,
		EditConfig: global.ApiEditConfig,
	})
}

// CleanCache 清除缓存
func CleanCache(r *ghttp.Request) {
	size, err := global.GfCache.Size(r.Context())
	x.FastResp(r, err, false).Resp()
	x.FastResp(r, global.GfCache.Clear(r.Context()), false).Resp()
	x.FastResp(r).SetData(g.Map{"size": size}).Resp()
}

// CacheSize 获取缓存大小
func CacheSize(r *ghttp.Request) {
	size, err := global.GfCache.Size(r.Context())
	x.FastResp(r, err, false).Resp()
	x.FastResp(r).SetData(g.Map{"size": size}).Resp()
}

// DelConfig 删除配置
func DelConfig(r *ghttp.Request) {
	x.FastResp(r, !global.ApiEditConfig, false).Resp("非法操作")
	serialNumberT := r.Get("serialNumber")
	x.FastResp(r, serialNumberT.IsEmpty(), false).Resp("参数错误")
	x.FastResp(r, global.ConfigDevice.DelConfig(serialNumberT.String())).Resp()
	x.FastResp(r, buffer.Buffer.Init(true), false).Resp("删除成功但重新初始化失败")
	x.FastResp(r).Resp()
}

// UpdateStatus 修改配置状态
func UpdateStatus(r *ghttp.Request) {
	x.FastResp(r, !global.ApiEditConfig, false).Resp("非法操作")
	serialNumberT := r.Get("serialNumber")
	statusT := r.Get("status")
	x.FastResp(r, serialNumberT.IsEmpty() || statusT.IsEmpty(), false).Resp("参数错误")
	x.FastResp(r, global.ConfigDevice.UpdateStatus(serialNumberT.String(), statusT.Int())).Resp()
	x.FastResp(r, buffer.Buffer.Init(true), false).Resp("修改成功但重新初始化失败")
	x.FastResp(r).Resp()
}
