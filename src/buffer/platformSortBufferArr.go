package buffer

import (
	"sort"
	"uniTranslate/src/types"
)

type PlatformSortBufferArr struct {
	buffer *BufferType
	level  [][]*types.TranslatePlatform
	idx    [][]int
}

func (r *PlatformSortBufferArr) Init(data *BufferType, platform string) BufferArrInterface {
	r.buffer = data
	level := r.buffer.GetLevel()

	// 预计算容量以减少内存重新分配
	totalPlatforms := 0
	for _, platforms := range level {
		totalPlatforms += len(platforms)
	}

	// 预分配内存
	platformArr := make([]*types.TranslatePlatform, 0, totalPlatforms)
	nowLevelArr := make([][]*types.TranslatePlatform, 0, len(level))

	// 一次遍历完成分类
	for _, platforms := range level {
		otherPlatforms := make([]*types.TranslatePlatform, 0, len(platforms))

		for _, tp := range platforms {
			// 创建新的 TranslatePlatform 对象
			newTP := &types.TranslatePlatform{
				Md5:      tp.Md5,
				Platform: tp.Platform,
				Status:   tp.Status,
				Level:    tp.Level,
				Cfg:      tp.Cfg,
				Type:     tp.Type,
			}

			if platform == tp.Type {
				platformArr = append(platformArr, newTP)
			} else {
				otherPlatforms = append(otherPlatforms, newTP)
			}
		}

		if len(otherPlatforms) > 0 {
			nowLevelArr = append(nowLevelArr, otherPlatforms)
		}
	}

	// 优化排序：只有当需要排序时才排序
	if len(platformArr) > 1 {
		sort.Slice(platformArr, func(i, j int) bool {
			return platformArr[i].Level < platformArr[j].Level
		})
	}

	// 使用预分配的切片来构建最终结果
	result := make([][]*types.TranslatePlatform, 0, len(nowLevelArr)+1)
	if len(platformArr) > 0 {
		result = append(result, platformArr)
	}
	result = append(result, nowLevelArr...)

	r.level = result
	r.idx = r.buffer.CreateIdxArr(result)
	return r
}

func (r *PlatformSortBufferArr) GetPlatformConfig(i0, i1 int) *types.TranslatePlatform {
	if i0 >= len(r.level) || i1 >= len(r.level[i0]) {
		return nil
	}
	return r.level[i0][i1]
}

func (r *PlatformSortBufferArr) GetIdx(i int) []int {
	if i >= len(r.idx) {
		return nil
	}
	return r.idx[i]
}
