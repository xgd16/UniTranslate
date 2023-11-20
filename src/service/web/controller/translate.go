package controller

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xgd16/gf-x-tool/x"
	"github.com/xgd16/gf-x-tool/xlib"
	"github.com/xgd16/gf-x-tool/xtranslate"
	"uniTranslate/src/global"
	queueHandler "uniTranslate/src/service/queue/handler"
	"uniTranslate/src/service/web/handler"
	"uniTranslate/src/types"
)

func Translate(r *ghttp.Request) {
	fromT := r.Get("from")
	toT := r.Get("to")
	textT := r.Get("text")
	platform := r.Get("platform").String()
	x.FastResp(r, fromT.IsEmpty() || toT.IsEmpty() || textT.IsEmpty(), false).Resp("参数错误")
	x.FastResp(r, platform != "" && !xlib.InArr(platform, []string{xtranslate.Baidu, xtranslate.Deepl, xtranslate.Google, xtranslate.YouDao}), false).Resp("不支持的平台")
	from := fromT.String()
	to := toT.String()
	text := textT.String()
	// 内容转换为md5
	md5 := gmd5.MustEncrypt(fmt.Sprintf("from:%s-to:%s-text:%s-platform:%s", from, to, text, platform))
	// 写入到缓存
	var (
		data *gvar.Var
		err  error
	)
	if global.CacheMode == "off" {
		var dataAny any
		dataAny, err = t(from, to, text, platform)
		data = gvar.New(dataAny)
	} else {
		data, err = global.GfCache.GetOrSetFunc(r.GetCtx(), fmt.Sprintf("Translate-%s", md5), func(ctx context.Context) (value any, err error) {
			return t(from, to, text, platform)
		}, 0)
	}
	queueHandler.RequestRecordQueue.Push(&types.RequestRecordData{
		ClientIp: r.GetClientIp(),
		Body:     r.GetBodyString(),
		Time:     gtime.Now().UnixMilli(),
		Ok:       err == nil,
		ErrMsg:   err,
	})
	x.FastResp(r, err, false).Resp("翻译失败请重试")
	x.FastResp(r).SetData(data.Map()).Resp()
}

func AddConfig(r *ghttp.Request) {
	t := new(types.TranslatePlatform)
	x.FastResp(r, r.GetStruct(t)).Resp()
	t.InitMd5()
	x.FastResp(r, t.Type != "" && !xlib.InArr(t.Type, []string{xtranslate.Baidu, xtranslate.Deepl, xtranslate.Google, xtranslate.YouDao}), false).Resp("不支持的平台")
	x.FastResp(r, !global.XDB.GetGJson().Get(fmt.Sprintf("xtranslate.%s", t.GetMd5())).IsEmpty(), false).Resp("已存在此配置")
	x.FastResp(r, global.XDB.Set("translate", t.GetMd5(), t), false).Resp("添加失败")
	x.FastResp(r, global.StatisticalProcess.CreateEvent(t)).Resp("添加失败")
	x.FastResp(r, global.Buffer.Init(), false).Resp("写入成功但重新初始化失败")
	x.FastResp(r).Resp()
}

func t(from, to, text, platform string) (value any, err error) {
	value, err = global.Buffer.Handler(from, to, text, platform, handler.Translate)
	queueHandler.CountRecordQueue.Push(&types.CountRecordData{
		Data: value.(*types.TranslateData),
		Ok:   err == nil,
	})
	return
}
