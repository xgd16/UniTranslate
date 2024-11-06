package controller

import (
	"sort"
	"uniTranslate/src/buffer"
	"uniTranslate/src/devices"
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

// AggregateTranslate 聚合翻译 (只支持单条翻译)
func AggregateTranslate(r *ghttp.Request) {
	req := new(types.AggregateTranslationReq)
	if err := r.Parse(req); err != nil {
		x.FastResp(r, true, false).Resp(err.Error())
	}
	result, err := logic.AggregateTranslate(r.Context(), req)
	x.FastResp(r, err, false).Resp("翻译失败请重试")
	x.FastResp(r).SetData(result).Resp()
}

// LibreTranslate 兼容实现
func LibreTranslate(r *ghttp.Request) {
	req := new(types.LibreTranslateReq)
	if err := r.Parse(req); err != nil {
		x.FastResp(r, true, false).Resp(err.Error())
	}
	data, err := logic.Translate(r.GetCtx(), r.GetClientIp(), &types.TranslateReq{
		From: req.Source,
		To:   req.Target,
		Text: []string{req.QueryStr},
	})
	if err != nil || len(data.Translate) <= 0 {
		r.Response.Status = 500
		r.Response.WriteJsonExit(g.Map{
			"error": "翻译失败请重试",
		})
	}
	dataItem := data.Translate[0]
	r.Response.WriteJsonExit(g.Map{
		"detectedLanguage": g.Map{
			"confidence": 0,
			"language":   dataItem.FromLang,
		},
		"translatedText": dataItem.Text,
	})
}

type IdeaTranslateVirtualQueryReq struct {
	Client string `json:"client"`
	Dj     string `json:"dj"`
	Dt     string `json:"dt"`
	Hl     string `json:"hl"`
	Ie     string `json:"ie"`
	Oe     string `json:"oe"`
	Sl     string `json:"sl"`
	Tk     string `json:"tk"`
	Tl     string `json:"tl"`
}

// GoogleSingleVirtual 谷歌翻译虚拟接口
func GoogleSingleVirtual(r *ghttp.Request) {
	if r.Get("key").String() != global.ServiceKey {
		r.Response.WriteStatusExit(404)
	}
	queryData := new(IdeaTranslateVirtualQueryReq)
	x.FastResp(r, r.GetQueryStruct(queryData), false).Resp("请求参数出错")
	q := r.Get("q").String()
	data, err := logic.Translate(r.GetCtx(), r.GetClientIp(), &types.TranslateReq{
		From: "auto",
		To:   queryData.Tl,
		Text: []string{q},
	})
	x.FastResp(r, err, false).Resp("翻译失败请重试")
	x.FastResp(r, len(data.Translate) <= 0, false).Resp("翻译失败")
	FromLang := data.Translate[0].FromLang
	r.Response.WriteJsonExit(g.Map{
		"sentences": g.Array{
			g.Map{
				"trans":   data.Translate[0].Text,
				"orig":    data.OriginalText[0],
				"backend": 3,
				"model_specification": g.Array{
					g.Map{},
				},
				"translation_engine_debug_info": g.Array{
					g.Map{
						"model_tracking": g.Map{
							"checkpoint_md5": "af64405095a399ceb1e05c7abb7cda66",
							"launch_doc":     "zh_en_2023q1.md",
						},
					},
				},
			},
			g.Map{
				"src_translit": "Zhīchí de píngtái 1",
			},
		},
		"src":        FromLang,
		"confidence": 1.0,
		"spell":      g.Map{},
		"ld_result": g.Map{
			"srclangs": g.Array{
				FromLang,
			},
			"srclangs_confidences": []float64{
				1.0,
			},
			"extended_srclangs": []string{
				FromLang,
			},
		},
	})
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
	data, err := logic.Translate(r.GetCtx(), r.GetClientIp(), req)
	// 处理结果
	x.FastResp(r, err, false).Resp("翻译失败请重试")
	x.FastResp(r).SetData(data).Resp()
}

// GetConfigList 获取配置列表
func GetConfigList(r *ghttp.Request) {
	// 获取配置驱动
	device, err := devices.GetConfigDevice()
	x.FastResp(r, err).Resp()
	// 获取配置
	config, err := device.GetConfig(true)
	x.FastResp(r, err).Resp()
	// 获取计数记录
	countRecordMap, err := devices.RecordHandler.GetCountRecord()
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

		configItem := g.Map{
			"id":          k,
			"level":       v.Level,
			"platform":    v.Platform,
			"status":      v.Status,
			"type":        v.Type,
			"countRecord": countRecord,
		}
		// 判断是否开启编辑配置
		if global.ApiEditConfig {
			configItem["cfg"] = v.Cfg
		}

		respData = append(respData, configItem)
	}
	// 按照level排序
	sort.Slice(respData, func(i, j int) bool {
		return gconv.Int(respData[i]["level"]) > gconv.Int(respData[j]["level"])
	})

	x.FastResp(r).SetData(respData).Resp()
}

// SaveConfig 添加配置
func SaveConfig(r *ghttp.Request) {
	t := new(types.TranslatePlatform)
	x.FastResp(r, r.GetStruct(t)).Resp()
	x.FastResp(r, t.Platform == "", false).Resp("名称不能为空")
	if t.Md5 == "" {
		t.InitMd5()
	}
	x.FastResp(r, !global.ApiEditConfig && t.Md5 != "", false).Resp("非法操作")
	x.FastResp(r, t.Type != "" && !xlib.InArr(t.Type, translate.TranslateModeList), false).Resp("不支持的平台")
	device, err := devices.GetConfigDevice()
	x.FastResp(r, err).Resp()
	_, ok, err := device.GetTranslateInfo(t.GetMd5())
	x.FastResp(r, err).Resp()
	x.FastResp(r, device.SaveConfig(t.GetMd5(), ok, t), false).Resp("添加失败")
	if !ok {
		x.FastResp(r, devices.RecordHandler.CreateEvent(t)).Resp("添加失败")
	}
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
	x.FastResp(r).SetData(global.LangJson).Resp()
}
