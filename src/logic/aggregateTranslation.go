package logic

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"golang.org/x/sync/errgroup"
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
		return nil, fmt.Errorf("获取配置驱动失败: %w", err)
	}

	config, err := device.GetConfig(false)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	// 处理平台列表
	req.Platform = getPlatformList(req.Platform, config)

	// 获取需要使用的配置
	needConfig := getNeedConfig(req.Platform, config)

	// 使用 mutex 保护 slice 并发写入
	var mu sync.Mutex
	translateResult = make([]*types.AggregateTranslateResult, 0, len(needConfig))

	// 使用 errgroup 替代 WaitGroup，更好地处理并发和错误
	eg := errgroup.Group{}

	for platform, configs := range needConfig {
		platform, configs := platform, configs // 创建局部变量避免闭包问题
		eg.Go(func() error {
			result := translateForPlatform(ctx, platform, configs, req, logger)

			mu.Lock()
			translateResult = append(translateResult, result)
			mu.Unlock()

			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		return translateResult, fmt.Errorf("聚合翻译过程出错: %w", err)
	}

	return translateResult, nil
}

// 获取平台列表
func getPlatformList(platforms []string, config map[string]*types.TranslatePlatform) []string {
	if len(platforms) > 0 {
		return platforms
	}

	platfromMap := make(map[string]struct{})
	for _, v := range config {
		platfromMap[v.Type] = struct{}{}
	}

	result := make([]string, 0, len(platfromMap))
	for platform := range platfromMap {
		result = append(result, platform)
	}
	return result
}

// 获取需要使用的配置
func getNeedConfig(platforms []string, config map[string]*types.TranslatePlatform) map[string][]*types.TranslatePlatform {
	needConfig := make(map[string][]*types.TranslatePlatform)

	for _, platform := range platforms {
		needConfig[platform] = make([]*types.TranslatePlatform, 0)
		for _, configItem := range config {
			if platform == configItem.Type {
				needConfig[platform] = append(needConfig[platform], configItem)
			}
		}
	}

	return needConfig
}

// 处理单个平台的翻译
func translateForPlatform(ctx context.Context, platform string, configs []*types.TranslatePlatform, req *types.AggregateTranslationReq, logger *glog.Logger) *types.AggregateTranslateResult {
	result := &types.AggregateTranslateResult{
		Platform: platform,
	}

	if len(configs) == 0 {
		result.ErrorStr = "不存在的翻译平台或没有配置"
		return result
	}

	config := configs[0]
	t, err := translate.GetTranslate(platform, config.Cfg)
	if err != nil {
		result.ErrorStr = "获取翻译器失败"
		logger.Errorf(ctx, "获取翻译器失败: %+v", err)
		return result
	}

	tResp, err := t.Translate(&translate.TranslateReq{
		From:     req.From,
		To:       req.To,
		Text:     []string{req.Text},
		Platform: platform,
	})

	if err != nil {
		result.ErrorStr = "翻译请求失败"
		logger.Errorf(ctx, "翻译请求失败: %+v", err)
		return result
	}

	if len(tResp) == 0 {
		result.ErrorStr = "翻译结果为空"
		return result
	}

	tRespItem := tResp[0]
	result.Translate = tRespItem.Text
	result.FromLang = tRespItem.FromLang

	return result
}
