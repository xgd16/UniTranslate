package logic

import (
	"context"
	"sync"
	"uniTranslate/src/devices"
	"uniTranslate/src/translate"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/frame/g"
)

func AggregateTranslate(ctx context.Context, req *types.AggregateTranslationReq) (translateResult []*types.AggregateTranslateResult, err error) {
	logger := g.Log()
	// 获取配置驱动
	device, err := devices.GetConfigDevice()
	if err != nil {
		return
	}
	config, err := device.GetConfig(false)
	if err != nil {
		return
	}
	// 判断是否没指定平台
	if len(req.Platform) <= 0 {
		platfromMap := make(map[string]struct{})
		for _, v := range config {
			platfromMap[v.Type] = struct{}{}
		}
		for platform := range platfromMap {
			req.Platform = append(req.Platform, platform)
		}
	}
	// 取出需要使用的配置
	needConfig := make(map[string][]*types.TranslatePlatform)
	for _, item := range req.Platform {
		for _, configItem := range config {
			if _, ok := needConfig[item]; !ok {
				needConfig[item] = make([]*types.TranslatePlatform, 0)
			}
			if item == configItem.Type {
				needConfig[item] = append(needConfig[item], configItem)
			}
		}
	}
	translateResult = make([]*types.AggregateTranslateResult, 0)
	wg := new(sync.WaitGroup)
	// 逐个平台翻译
	for platform, configs := range needConfig {
		wg.Add(1)
		go func(platform string, configs []*types.TranslatePlatform) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(ctx, "聚合翻译出错", err)
				}
				wg.Done()
			}()
			result := new(types.AggregateTranslateResult)
			result.Platform = platform
			if len(configs) <= 0 {
				result.ErrorStr = "不存在的翻译平台或没有配置"
			} else {
				config := configs[0]
				// 获取翻译器
				t, err := translate.GetTranslate(platform, config.Cfg)
				if err != nil {
					result.ErrorStr = "翻译出错请稍后再试"
					logger.Error(ctx, err)
				} else {
					// 执行翻译
					tResp, err := t.Translate(&translate.TranslateReq{
						From:     req.From,
						To:       req.To,
						Text:     []string{req.Text},
						Platfrom: platform,
					})
					if err != nil {
						result.ErrorStr = "翻译出错请稍后再试"
						logger.Error(ctx, err)
					} else {
						// 判断是否有翻译结果
						if len(tResp) <= 0 {
							result.ErrorStr = "翻译失败"
						} else {
							tRespItem := tResp[0]
							result.Translate = tRespItem.Text
							result.FromLang = tRespItem.FromLang
						}
					}
				}
			}
			translateResult = append(translateResult, result)
		}(platform, configs)
	}
	wg.Wait()
	return
}
