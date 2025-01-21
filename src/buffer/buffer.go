package buffer

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"sync"
	"uniTranslate/src/devices"
	"uniTranslate/src/global"
	queueHandler "uniTranslate/src/service/queue/handler"
	"uniTranslate/src/translate"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xgd16/gf-x-tool/xlib"
	"github.com/xgd16/gf-x-tool/xmonitor"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var Buffer = new(BufferType)

type BufferType struct {
	m        sync.Mutex
	level    [][]*types.TranslatePlatform
	levelArr *garray.IntArray
	num      int
	idx      [][]int
}

func (t *BufferType) GetLevel() [][]*types.TranslatePlatform {
	return t.level
}

func (t *BufferType) GetIdx() [][]int {
	return t.idx
}

func (t *BufferType) Handler(req *translate.TranslateReq, fn func(config *types.TranslatePlatform, req *translate.TranslateReq) (*types.TranslateData, error)) (s *types.TranslateData, e error) {
	// 获取 buffer array 的操作应该在锁外完成
	var bufferArr BufferArrInterface
	if req.Platform == "" {
		bufferArr = new(RandomSortBufferArr)
	} else {
		bufferArr = new(PlatformSortBufferArr)
	}

	t.m.Lock()
	bufferArr.Init(t, req.Platform)
	t.m.Unlock()

	ctx := gctx.New()

	for i := 0; i < t.num; i++ {
		var p *types.TranslatePlatform
		func() {
			t.m.Lock()
			defer t.m.Unlock()

			idx := bufferArr.GetIdx(i)
			p = bufferArr.GetPlatformConfig(idx[0], idx[1])
			if p.Status == 0 {
				return
			}
		}()

		if p == nil || p.Status == 0 {
			continue
		}

		translateResp, err := fn(p, req)

		// 处理监控指标
		t.handleMetrics(err, p, translateResp, req)

		if err != nil {
			e = fmt.Errorf("调用翻译失败: %w", err)
			t.recordError(ctx, req, err, p)
			continue
		}

		translateResp.Md5 = p.Md5
		return translateResp, nil
	}

	if e == nil {
		e = errors.New("所有翻译平台均调用失败")
	}
	return
}

// 新增的辅助方法，用于处理监控指标
func (t *BufferType) handleMetrics(err error, p *types.TranslatePlatform, resp *types.TranslateData, req *translate.TranslateReq) {
	if global.RunMode != global.HttpMode {
		return
	}

	status := xlib.IF(err == nil, "success", "error")
	xmonitor.MetricHttpRequestTotal.WithLabelValues(fmt.Sprintf("%s_%s", status, p.Platform)).Inc()

	if err == nil && resp != nil {
		fontCount := gconv.Float64(gstr.LenRune(gstr.Join(resp.OriginalText, "")))
		xmonitor.MetricHttpRequestTotal.WithLabelValues(fmt.Sprintf("fontCount_%s", p.Platform)).Add(fontCount)
	}
}

// 新增的辅助方法，用于记录错误
func (t *BufferType) recordError(ctx context.Context, req *translate.TranslateReq, err error, p *types.TranslatePlatform) {
	queueHandler.RequestRecordQueue.Push(&types.RequestRecordData{
		ClientIp: req.HttpReq.ClientIp,
		Body: &types.TranslateReq{
			From:     req.From,
			To:       req.To,
			Text:     req.Text,
			Platform: req.Platform,
		},
		Time:     gtime.Now().UnixMilli(),
		Ok:       false,
		ErrMsg:   err,
		Platform: fmt.Sprintf("%s [ %s ]", p.Type, p.Platform),
		TraceId:  gtrace.GetTraceID(req.HttpReq.Context),
	})

	queueHandler.CountRecordQueue.Push(&types.CountRecordData{
		Data: &types.TranslateData{
			Md5: p.Md5,
		},
		Ok: false,
	})

	g.Log().Error(ctx, err)
}

func (t *BufferType) Init(refresh bool) (err error) {
	t.m.Lock()
	defer t.m.Unlock()
	// 初始化数据
	device, err := devices.GetConfigDevice()
	if err != nil {
		return
	}
	config, err := device.GetConfig(refresh)
	if err != nil {
		return
	}
	t.level, t.idx, t.num = t.getLevelSort(config)
	return
}

func (t *BufferType) getLevelSort(data map[string]*types.TranslatePlatform) (arr [][]*types.TranslatePlatform, idxArr [][]int, num int) {
	// 格式化结构以及统计个数
	num = 0
	levelArrT := garray.NewIntArray()
	for k, platform := range data {
		num += 1
		platform.Md5 = k
		levelArrT.PushRight(platform.Level)
	}
	t.levelArr = levelArrT.Unique().Sort()
	l := t.levelArr.Len()
	// 创建用于操作的结构数据
	arr = make([][]*types.TranslatePlatform, l)
	for _, platform := range data {
		idx := t.levelArr.Search(platform.Level)
		arr[idx] = append(arr[idx], platform)
	}
	// 创建横向索引
	idxArr = t.CreateIdxArr(arr)
	return
}

func (t *BufferType) CreateIdxArr(arr [][]*types.TranslatePlatform) [][]int {
	idxArr := make([][]int, 0)
	for k, v := range arr {
		for k1 := range v {
			idxArr = append(idxArr, []int{k, k1})
		}
	}
	return idxArr
}

type BufferArrInterface interface {
	// Init 初始化翻译平台对象
	Init(data *BufferType, platform string) BufferArrInterface
	// GetPlatformConfig 获取平台配置
	GetPlatformConfig(i0, i1 int) *types.TranslatePlatform
	// GetIdx 获取 index
	GetIdx(i int) []int
}
