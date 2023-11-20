package types

import (
	"errors"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/xgd16/gf-x-tool/xstorage"
	"sync"
)

type BufferType struct {
	m        sync.Mutex
	level    [][]*TranslatePlatform
	levelArr *garray.IntArray
	num      int
	idx      [][]int
	xdb      *xstorage.XDB
}

func (t *BufferType) GetLevel() [][]*TranslatePlatform {
	return t.level
}

func (t *BufferType) GetIdx() [][]int {
	return t.idx
}

func (t *BufferType) SetXDB(xdb *xstorage.XDB) *BufferType {
	t.xdb = xdb
	return t
}

func (t *BufferType) Handler(from, to, text, platform string, fn func(*TranslatePlatform, string, string, string) (*TranslateData, error)) (s *TranslateData, e error) {
	t.m.Lock()
	var bufferArr BufferArrInterface
	if platform == "" {
		bufferArr = new(RandomSortBufferArr)
	} else {
		bufferArr = new(PlatformSortBufferArr)
	}
	bufferArr.Init(t, platform)
	t.m.Unlock()
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
			g.Log().Error(gctx.New(), "调用翻译失败", err)
			continue
		}
		t.Md5 = p.md5
		return t, nil
	}
	e = errors.New("翻译失败")
	return
}

func (t *BufferType) Init() (err error) {
	t.m.Lock()
	defer t.m.Unlock()
	// 初始化数据
	data := new(map[string]*TranslatePlatform)
	err = t.xdb.GetGJson().Get("translate").Scan(data)
	if err != nil {
		return
	}
	t.level, t.idx, t.num = t.getLevelSort(*data)
	return
}

func (t *BufferType) getLevelSort(data map[string]*TranslatePlatform) (arr [][]*TranslatePlatform, idxArr [][]int, num int) {
	// 格式化结构以及统计个数
	num = 0
	levelArrT := garray.NewIntArray()
	for k, platform := range data {
		num += 1
		platform.md5 = k
		levelArrT.PushRight(platform.Level)
	}
	t.levelArr = levelArrT.Unique().Sort()
	l := t.levelArr.Len()
	// 创建用于操作的结构数据
	arr = make([][]*TranslatePlatform, l)
	for _, platform := range data {
		idx := t.levelArr.Search(platform.Level)
		arr[idx] = append(arr[idx], platform)
	}
	// 创建横向索引
	idxArr = t.CreateIdxArr(arr)
	return
}

func (t *BufferType) CreateIdxArr(arr [][]*TranslatePlatform) [][]int {
	idxArr := make([][]int, 0)
	for k, v := range arr {
		for k1, _ := range v {
			idxArr = append(idxArr, []int{k, k1})
		}
	}
	return idxArr
}

type BufferArrInterface interface {
	// Init 初始化翻译平台对象
	Init(data *BufferType, platform string) BufferArrInterface
	// GetPlatformConfig 获取平台配置
	GetPlatformConfig(i0, i1 int) *TranslatePlatform
	// GetIdx 获取 index
	GetIdx(i int) []int
}
