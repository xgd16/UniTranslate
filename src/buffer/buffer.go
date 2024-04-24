package buffer

import (
	"errors"
	"fmt"
	"sync"
	"uniTranslate/src/global"
	queueHandler "uniTranslate/src/service/queue/handler"
	"uniTranslate/src/types"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"

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

func (t *BufferType) Handler(r *ghttp.Request, from, to, text, platform string, fn func(*types.TranslatePlatform, string, string, string) (*types.TranslateData, error)) (s *types.TranslateData, e error) {
	t.m.Lock()
	var bufferArr BufferArrInterface
	if platform == "" {
		bufferArr = new(RandomSortBufferArr)
	} else {
		bufferArr = new(PlatformSortBufferArr)
	}
	bufferArr.Init(t, platform)
	t.m.Unlock()
	// 创建上下文
	ctx := gctx.New()
	// 循环处理数据
	for i := 0; i < t.num; i++ {
		t.m.Lock()
		// 获取操作对象
		idx := bufferArr.GetIdx(i)
		p := bufferArr.GetPlatformConfig(idx[0], idx[1])
		// 释放锁
		t.m.Unlock()
		// 调用处理
		t, err := fn(p, from, to, text)
		if err != nil {
			e = fmt.Errorf("调用翻译失败 %s", err)
			queueHandler.RequestRecordQueue.Push(&types.RequestRecordData{
				ClientIp: r.GetClientIp(),
				Body:     gstr.TrimAll(r.GetBodyString()),
				Time:     gtime.Now().UnixMilli(),
				Ok:       err == nil,
				ErrMsg:   err,
				Platform: fmt.Sprintf("%s [ %s ]", p.Type, p.Platform),
				TraceId:  gtrace.GetTraceID(r.Context()),
			})
			g.Log().Error(ctx, e)
			continue
		}
		t.Md5 = p.Md5
		return t, nil
	}
	if e == nil {
		e = errors.New("翻译失败")
	}
	return
}

func (t *BufferType) Init(refresh bool) (err error) {
	t.m.Lock()
	defer t.m.Unlock()
	// 初始化数据
	device, err := global.GetConfigDevice()
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
		if platform.Status == 0 {
			continue
		}
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
