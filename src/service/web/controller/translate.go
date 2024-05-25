package controller

import (
	"sort"
	"uniTranslate/src/buffer"
	"uniTranslate/src/global"
	"uniTranslate/src/logic"
	"uniTranslate/src/translate"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xgd16/gf-x-tool/x"
	"github.com/xgd16/gf-x-tool/xlib"
)

// AggregateTranslation 聚合翻译 (只支持单条翻译)
func AggregateTranslate(r *ghttp.Request) {
	req := new(types.AggregateTranslationReq)
	if err := r.Parse(req); err != nil {
		x.FastResp(r, true, false).Resp(err.Error())
	}
	result, err := logic.AggregateTranslate(r.Context(), req)
	x.FastResp(r, err, false).Resp("翻译失败请重试")
	x.FastResp(r).SetData(result).Resp()
}

// Translate 翻译
func Translate(r *ghttp.Request) {
	req := new(types.TranslateReq)
	if err := r.Parse(req); err != nil {
		x.FastResp(r, true, false).Resp(err.Error())
	}
	x.FastResp(r, req.Platform != "" && !xlib.InArr(req.Platform, translate.TranslateModeList), false).Resp("不支持的平台")
	x.FastResp(r, req.To == "auto", false).Resp("转换后语言不支持 auto")
	// 翻译
	data, err := logic.Translate(r, req)
	// 处理结果
	x.FastResp(r, err, false).Resp("翻译失败请重试")
	x.FastResp(r).SetData(data).Resp()
}

// GetConfigList 获取配置列表
func GetConfigList(r *ghttp.Request) {
	// 获取配置驱动
	device, err := global.GetConfigDevice()
	x.FastResp(r, err).Resp()
	// 获取配置
	config, err := device.GetConfig(true)
	x.FastResp(r, err).Resp()
	// 获取计数记录
	countRecordMap, err := global.StatisticalProcess.GetCountRecord()
	x.FastResp(r, err, false).Resp()

	respData := make([]map[string]any, 0)
	for k, v := range config {
		var (
			countRecord *types.CountRecord
			ok          bool
		)
		if countRecord, ok = countRecordMap[k]; !ok {
			countRecord = new(types.CountRecord)
		}
		respData = append(respData, g.Map{
			"id":          k,
			"level":       v.Level,
			"platform":    v.Platform,
			"status":      v.Status,
			"type":        v.Type,
			"countRecord": countRecord,
		})
	}
	// 按照level排序
	sort.Slice(respData, func(i, j int) bool {
		return gconv.Int(respData[i]["level"]) > gconv.Int(respData[j]["level"])
	})

	x.FastResp(r).SetData(respData).Resp()
}

// AddConfig 添加配置
func AddConfig(r *ghttp.Request) {
	t := new(types.TranslatePlatform)
	x.FastResp(r, r.GetStruct(t)).Resp()
	x.FastResp(r, t.Platform == "", false).Resp("名称不能为空")
	t.InitMd5()
	x.FastResp(r, t.Type != "" && !xlib.InArr(t.Type, translate.TranslateModeList), false).Resp("不支持的平台")
	device, err := global.GetConfigDevice()
	x.FastResp(r, err).Resp()
	_, ok, err := device.GetTranslateInfo(t.GetMd5())
	x.FastResp(r, err).Resp()
	x.FastResp(r, ok, false).Resp("已存在此配置")
	x.FastResp(r, device.SaveConfig(t.GetMd5(), t), false).Resp("添加失败")
	x.FastResp(r, global.StatisticalProcess.CreateEvent(t)).Resp("添加失败")
	x.FastResp(r, buffer.Buffer.Init(true), false).Resp("写入成功但重新初始化失败")
	x.FastResp(r).Resp()
}

// RefreshConfigCache 刷新配置缓存
func RefreshConfigCache(r *ghttp.Request) {
	x.FastResp(r, buffer.Buffer.Init(true), false).Resp("写入成功但重新初始化失败")
	x.FastResp(r).Resp()
}

// GetLangList 获取语言列表
func GetLangList(r *ghttp.Request) {
	x.FastResp(r).SetData(translate.BaseTranslateConf[translate.YouDaoTranslateMode]).Resp()
}
