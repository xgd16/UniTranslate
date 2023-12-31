package buffer

import (
	"github.com/gogf/gf/v2/container/garray"
	"uniTranslate/src/types"
)

type RandomSortBufferArr struct {
	buffer *BufferType
	level  [][]*types.TranslatePlatform
	idx    [][]int
}

func (r *RandomSortBufferArr) Init(data *BufferType, _ string) BufferArrInterface {
	r.buffer = data
	level := r.buffer.GetLevel()
	// 随机排序当前序列
	nowLevelArr := make([][]*types.TranslatePlatform, len(level))
	for i, platforms := range level {
		// 数据不大于1的跳过
		if len(platforms) <= 1 {
			nowLevelArr[i] = platforms
			continue
		}
		// 打乱组中数据顺序
		nowLevelArr[i] = func(platforms []*types.TranslatePlatform) []*types.TranslatePlatform {
			var pArr []*types.TranslatePlatform
			l := len(platforms)
			a := garray.NewIntArray()
			for n := 0; n < l; n++ {
				a.PushRight(n)
			}
			for _, v := range a.Shuffle().Slice() {
				pArr = append(pArr, platforms[v])
			}
			return pArr
		}(platforms)
	}
	r.level = nowLevelArr
	r.idx = r.buffer.GetIdx()
	return r
}

func (r *RandomSortBufferArr) GetPlatformConfig(i0, i1 int) *types.TranslatePlatform {
	return r.level[i0][i1]
}

func (r *RandomSortBufferArr) GetIdx(i int) []int {
	return r.idx[i]
}
