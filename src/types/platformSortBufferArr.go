package types

import (
	"sort"
)

type PlatformSortBufferArr struct {
	buffer *BufferType
	level  [][]*TranslatePlatform
	idx    [][]int
}

func (r *PlatformSortBufferArr) Init(data *BufferType, platform string) BufferArrInterface {
	r.buffer = data
	level := deepCopy2DArray(r.buffer.GetLevel())
	// 从等级排列重新按照平台建立新结构
	nowLevelArr := make([][]*TranslatePlatform, 0)
	platformArr := make([]*TranslatePlatform, 0)
	for _, platforms := range level {
		arr := make([]*TranslatePlatform, 0)
		for _, translatePlatform := range platforms {
			if platform == translatePlatform.Type {
				platformArr = append(platformArr, translatePlatform)
			} else {
				arr = append(arr, translatePlatform)
			}
		}
		if len(arr) > 0 {
			nowLevelArr = append(nowLevelArr, arr)
		}
	}
	// 按照翻译优先级从小到大排序
	sort.Slice(platformArr, func(i, j int) bool {
		return platformArr[i].Level < platformArr[j].Level
	})
	nowLevelArr = append([][]*TranslatePlatform{platformArr}, nowLevelArr...)
	r.level = nowLevelArr
	r.idx = r.buffer.CreateIdxArr(nowLevelArr)
	return r
}

func (r *PlatformSortBufferArr) GetPlatformConfig(i0, i1 int) *TranslatePlatform {
	return r.level[i0][i1]
}

func (r *PlatformSortBufferArr) GetIdx(i int) []int {
	return r.idx[i]
}

func deepCopy2DArray(input [][]*TranslatePlatform) [][]*TranslatePlatform {
	if input == nil {
		return nil
	}
	// 创建一个新的二维切片
	result := make([][]*TranslatePlatform, len(input))
	// 遍历原始数组的每个元素
	for i, row := range input {
		// 创建一个新的 TranslatePlatform 切片
		newRow := make([]*TranslatePlatform, len(row))
		// 遍历原始数组元素的每个元素，并复制它们
		for j, element := range row {
			if element != nil {
				// 这里复制 TranslatePlatform 的字段值
				// 例如：newElement.Field1 = element.Field1
				// ...
				newRow[j] = &TranslatePlatform{
					element.md5,
					element.Platform,
					element.Status,
					element.Level,
					element.Cfg,
					element.Type,
				}
			}
		}
		// 将新的行添加到结果中
		result[i] = newRow
	}
	return result
}
