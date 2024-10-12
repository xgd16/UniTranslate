package buffer

import (
	"uniTranslate/src/lib"
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
		newPlatforms := make([]*types.TranslatePlatform, len(platforms))
		copy(newPlatforms, platforms)
		// 打乱组中数据顺序
		lib.Shuffle(newPlatforms)
		nowLevelArr[i] = newPlatforms
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
