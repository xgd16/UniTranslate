package logic

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"uniTranslate/src/buffer"
	"uniTranslate/src/global"
	queueHandler "uniTranslate/src/service/queue/handler"
	"uniTranslate/src/service/web/handler"
	"uniTranslate/src/translate"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/xgd16/gf-x-tool/xmonitor"
)

// Translate 翻译
func Translate(ctx context.Context, ip string, req *types.TranslateReq) (data *types.TranslateData, err error) {
	textStr := gstr.Join(req.Text, "\n")
	// 内容转换为md5
	var keyStr string
	if global.ServerConfig.CachePlatform {
		keyStr = fmt.Sprintf("to:%s-text:%s-platform:%s", req.To, textStr, req.Platform)
	} else {
		keyStr = fmt.Sprintf("to:%s-text:%s", req.To, textStr)
	}
	cacheId := gmd5.MustEncrypt(keyStr)
	// 记录从翻译到获取到结果的时间
	startTime := gtime.Now().UnixMilli()
	// 创建所需要的参数
	translateReq := &translate.TranslateReq{
		HttpReq: &translate.TranslateHttpReq{
			ClientIp: ip,
			Context:  ctx,
		},
		From:     req.From,
		To:       req.To,
		Platform: req.Platform,
		Text:     req.Text,
		TextStr:  textStr,
	}
	isCache := true
	// 判断是否进行缓存
	if global.ServerConfig.CacheMode == "off" {
		isCache = false
		data, err = translateHandler(cacheId, translateReq)
	} else {
		data, err = hotColdDataTranslateHandler(ctx, cacheId, translateReq, &isCache)
		if err != nil {
			return
		}
	}
	// 统计命中缓存次数
	if isCache {
		xmonitor.MetricHttpRequestTotal.WithLabelValues("cache_translate_count").Inc()
	}
	nowTime := gtime.Now().UnixMilli()
	// 记录翻译
	queueHandler.RequestRecordQueue.Push(&types.RequestRecordData{
		ClientIp: ip,
		Body:     req,
		CacheId:  cacheId,
		Time:     nowTime,
		Ok:       err == nil,
		ErrMsg:   err,
		Platform: data.Platform,
		TakeTime: int(nowTime - startTime),
		TraceId:  gtrace.GetTraceID(ctx),
	})
	return
}

// translateHandler 翻译处理
func translateHandler(cacheId string, req *translate.TranslateReq) (data *types.TranslateData, err error) {
	data, err = buffer.Buffer.Handler(req, handler.Translate)
	if data != nil {
		data.CacheId = cacheId
		// 缓存写入数据库
		if global.ServerConfig.CacheWriteToStorage {
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

func hotColdDataTranslateHandler(ctx context.Context, cacheId string, req *translate.TranslateReq, isCache *bool) (data *types.TranslateData, err error) {
	// 获取连接
	rc, _ := g.Redis().Conn(ctx)
	// 释放连接
	defer func(rc gredis.Conn, ctx context.Context) {
		_ = rc.Close(ctx)
	}(rc, ctx)

	var res *gvar.Var
	// 缓存key
	cacheKey := "Translate:" + cacheId
	// 热数据
	res, err = rc.Do(ctx, "Get", cacheKey)
	if err != nil {
		return nil, err
	}
	if !res.IsEmpty() {
		if err = res.Scan(&data); err != nil {
			return nil, err
		}
		return
	}
	// 冷数据
	err = g.Model("translate_cache").Where("cacheId", cacheId).Scan(&data)
	if err != nil {
		return nil, err
	}
	if data == nil {
		*isCache = false
		data, err = translateHandler(cacheId, req)
		if err != nil {
			return nil, err
		}
	}
	_, err = rc.Do(ctx, "Set", cacheKey, interface{}(data))
	if err != nil {
		return nil, err
	}

	return
}
