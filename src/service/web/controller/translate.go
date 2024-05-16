package controller

import (
	"context"
	"fmt"
	"sort"
	"uniTranslate/src/buffer"
	"uniTranslate/src/global"
	"uniTranslate/src/logic"
	queueHandler "uniTranslate/src/service/queue/handler"
	"uniTranslate/src/service/web/handler"
	"uniTranslate/src/translate"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
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
	fromT := r.Get("from")
	toT := r.Get("to")
	textT := r.Get("text")
	platform := r.Get("platform").String()
	x.FastResp(r, fromT.IsEmpty() || toT.IsEmpty() || textT.IsEmpty(), false).Resp("参数错误")
	x.FastResp(r, platform != "" && !xlib.InArr(platform, translate.TranslateModeList), false).Resp("不支持的平台")
	from := fromT.String()
	to := toT.String()
	x.FastResp(r, to == "auto", false).Resp("转换后语言不支持 auto")
	text := textT.Strings()
	textStr := gstr.Join(text, "\n")
	// 内容转换为md5
	var keyStr string
	if global.CachePlatform {
		keyStr = fmt.Sprintf("to:%s-text:%s-platform:%s", to, textStr, platform)
	} else {
		keyStr = fmt.Sprintf("to:%s-text:%s", to, textStr)
	}
	md5 := gmd5.MustEncrypt(keyStr)
	// 写入到缓存
	var (
		data *gvar.Var
		err  error
	)
	// 记录从翻译到获取到结果的时间
	startTime := gtime.Now().UnixMilli()
	req := &translate.TranslateReq{
		From:     from,
		To:       to,
		Platfrom: platform,
		Text:     text,
		TextStr:  textStr,
	}
	// 判断是否进行缓存
	if global.CacheMode == "off" {
		var dataAny any
		dataAny, err = t(r, req)
		data = gvar.New(dataAny)
	} else {
		data, err = global.GfCache.GetOrSetFunc(r.GetCtx(), fmt.Sprintf("Translate:%s", md5), func(ctx context.Context) (value any, err error) {
			return t(r, req)
		}, 0)
	}
	endTime := gtime.Now().UnixMilli()
	// 转换为map
	dataMap := data.MapStrVar()
	// 记录翻译
	queueHandler.RequestRecordQueue.Push(&types.RequestRecordData{
		ClientIp: r.GetClientIp(),
		Body:     gstr.Trim(r.GetBodyString()),
		Time:     gtime.Now().UnixMilli(),
		Ok:       err == nil,
		ErrMsg:   err,
		Platform: dataMap["platform"].String(),
		TakeTime: int(endTime - startTime), // 获取到获取翻译的毫秒数
		TraceId:  gtrace.GetTraceID(r.Context()),
	})
	x.FastResp(r, err, false).Resp("翻译失败请重试")
	x.FastResp(r).SetData(dataMap).Resp()
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

func t(r *ghttp.Request, req *translate.TranslateReq) (value any, err error) {
	var data *types.TranslateData
	data, err = buffer.Buffer.Handler(r, req, handler.Translate)
	value = data

	if data != nil {
		// 缓存写入数据库
		if global.CacheWriteToStorage {
			queueHandler.SaveQueue.Push(&types.SaveData{
				Data: data,
			})
		}
		// 翻译计数
		queueHandler.CountRecordQueue.Push(&types.CountRecordData{
			Data: data,
			Ok:   err == nil,
		})
	}
	return
}

// GetLangList 获取语言列表
func GetLangList(r *ghttp.Request) {
	x.FastResp(r).SetData(translate.BaseTranslateConf[translate.YouDaoTranslateMode]).Resp()
}
