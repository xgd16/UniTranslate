package controller

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/text/gstr"
	"uniTranslate/src/buffer"
	"uniTranslate/src/global"
	"uniTranslate/src/lib"
	queueHandler "uniTranslate/src/service/queue/handler"
	"uniTranslate/src/service/web/handler"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xgd16/gf-x-tool/x"
	"github.com/xgd16/gf-x-tool/xlib"
	"github.com/xgd16/gf-x-tool/xtranslate"
)

var translateModeList = []string{
	xtranslate.Baidu,
	xtranslate.Deepl,
	xtranslate.Google,
	xtranslate.YouDao,
	lib.ChatGptTranslateMode,
	lib.XunFeiTranslateMode,
	lib.XunFeiNiuTranslateMode,
}

func Translate(r *ghttp.Request) {
	fromT := r.Get("from")
	toT := r.Get("to")
	textT := r.Get("text")
	platform := r.Get("platform").String()
	x.FastResp(r, fromT.IsEmpty() || toT.IsEmpty() || textT.IsEmpty(), false).Resp("参数错误")
	x.FastResp(r, platform != "" && !xlib.InArr(platform, translateModeList), false).Resp("不支持的平台")
	from := fromT.String()
	to := toT.String()
	text := textT.String()
	// 内容转换为md5
	var md5 string
	if global.CachePlatform {
		md5 = gmd5.MustEncrypt(fmt.Sprintf("to:%s-text:%s-platform:%s", to, text, platform))
	} else {
		md5 = gmd5.MustEncrypt(fmt.Sprintf("to:%s-text:%s", to, text))
	}
	// 写入到缓存
	var (
		data *gvar.Var
		err  error
	)
	// 记录从翻译到获取到结果的时间
	startTime := gtime.Now().UnixMilli()
	// 判断是否进行缓存
	if global.CacheMode == "off" {
		var dataAny any
		dataAny, err = t(r, from, to, text, platform)
		data = gvar.New(dataAny)
	} else {
		data, err = global.GfCache.GetOrSetFunc(r.GetCtx(), fmt.Sprintf("Translate-%s", md5), func(ctx context.Context) (value any, err error) {
			return t(r, from, to, text, platform)
		}, 0)
	}
	endTime := gtime.Now().UnixMilli()
	// 转换为map
	dataMap := data.MapStrVar()
	// 记录翻译
	queueHandler.RequestRecordQueue.Push(&types.RequestRecordData{
		ClientIp: r.GetClientIp(),
		Body:     gstr.TrimAll(r.GetBodyString()),
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

func GetConfigList(r *ghttp.Request) {
	dataMap := global.XDB.GetGJson().Get("translate").MapStrVar()
	respData := make([]map[string]any, 0)
	for k, v := range dataMap {
		item := v.MapStrAny()
		respData = append(respData, g.Map{
			"id":       k,
			"level":    item["level"],
			"platform": item["platform"],
			"status":   item["status"],
			"type":     item["type"],
		})
	}
	x.FastResp(r).SetData(respData).Resp()
}

func AddConfig(r *ghttp.Request) {
	t := new(types.TranslatePlatform)
	x.FastResp(r, r.GetStruct(t)).Resp()
	x.FastResp(r, t.Platform == "", false).Resp("名称不能为空")
	t.InitMd5()
	x.FastResp(r, t.Type != "" && !xlib.InArr(t.Type, translateModeList), false).Resp("不支持的平台")
	x.FastResp(r, !global.XDB.GetGJson().Get(fmt.Sprintf("xtranslate.%s", t.GetMd5())).IsEmpty(), false).Resp("已存在此配置")
	x.FastResp(r, global.XDB.Set("translate", t.GetMd5(), t), false).Resp("添加失败")
	x.FastResp(r, global.StatisticalProcess.CreateEvent(t)).Resp("添加失败")
	x.FastResp(r, buffer.Buffer.Init(), false).Resp("写入成功但重新初始化失败")
	x.FastResp(r).Resp()
}

func t(r *ghttp.Request, from, to, text, platform string) (value any, err error) {
	var data *types.TranslateData
	data, err = buffer.Buffer.Handler(r, from, to, text, platform, handler.Translate)
	value = data
	// 触发写入
	queueHandler.SaveQueue.Push(&types.SaveData{
		Data: data,
	})
	// 写入到缓存
	queueHandler.CountRecordQueue.Push(&types.CountRecordData{
		Data: data,
		Ok:   err == nil,
	})
	return
}
