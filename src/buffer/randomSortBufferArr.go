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

// Init 初始化并随机排序
func (r *RandomSortBufferArr) Init(data *BufferType, _ string) BufferArrInterface {
	r.buffer = data
	level := r.buffer.GetLevel()

	// 预分配内存
	r.level = make([][]*types.TranslatePlatform, len(level))

	// 随机排序当前序列
	for i, platforms := range level {
		// 数据不大于1的直接赋值
		if len(platforms) <= 1 {
			r.level[i] = platforms
			continue
		}

		// 直接在原切片上进行随机排序
		r.level[i] = make([]*types.TranslatePlatform, len(platforms))
		copy(r.level[i], platforms)
		lib.Shuffle(r.level[i])
	}

	r.idx = r.buffer.GetIdx()
	return r
}

// GetPlatformConfig 获取平台配置
func (r *RandomSortBufferArr) GetPlatformConfig(i0, i1 int) *types.TranslatePlatform {
	if i0 >= len(r.level) || i1 >= len(r.level[i0]) {
		return nil
	}
	return r.level[i0][i1]
}

// GetIdx 获取索引
func (r *RandomSortBufferArr) GetIdx(i int) []int {
	if i >= len(r.idx) {
		return nil
	}
	return r.idx[i]
}
