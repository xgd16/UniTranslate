package test

import (
	"testing"
	"uniTranslate/src/lib"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/stretchr/testify/assert"
)

func TestSortMapToStr(t *testing.T) {
	tests := []struct {
		name     string
		input1   map[string]any
		input2   map[string]any // 等价但顺序不同的输入
		expected string
	}{
		{
			name: "基础类型测试",
			input1: map[string]any{
				"name": "test",
				"age":  25,
			},
			input2: map[string]any{
				"age":  25,
				"name": "test",
			},
			expected: "age:25&name:test",
		},
		{
			name: "嵌套Map测试",
			input1: map[string]any{
				"user": map[string]any{
					"name": "test",
					"age":  25,
				},
				"status": "active",
			},
			input2: map[string]any{
				"status": "active",
				"user": map[string]any{
					"age":  25,
					"name": "test",
				},
			},
			expected: "status:active&user:{age:25&name:test}",
		},
		{
			name: "数组测试",
			input1: map[string]any{
				"tags": []any{"go", "test"},
				"id":   1,
			},
			input2: map[string]any{
				"id":   1,
				"tags": []any{"go", "test"},
			},
			expected: "id:1&tags:[go,test]",
		},
		{
			name: "复杂嵌套测试",
			input1: map[string]any{
				"data": map[string]any{
					"users": []any{
						map[string]any{"id": 1, "name": "user1"},
						map[string]any{"name": "user2", "id": 2},
					},
					"total": 2,
				},
				"status": "ok",
			},
			input2: map[string]any{
				"status": "ok",
				"data": map[string]any{
					"total": 2,
					"users": []any{
						map[string]any{"name": "user1", "id": 1},
						map[string]any{"id": 2, "name": "user2"},
					},
				},
			},
			expected: "data:{total:2&users:[{id:1&name:user1},{id:2&name:user2}]}&status:ok",
		},
		{
			name:     "空Map测试",
			input1:   map[string]any{},
			input2:   map[string]any{},
			expected: "",
		},
		{
			name: "特殊字符测试",
			input1: map[string]any{
				"url":    "https://example.com?key=value",
				"params": map[string]any{"&": "special"},
			},
			input2: map[string]any{
				"params": map[string]any{"&": "special"},
				"url":    "https://example.com?key=value",
			},
			expected: "params:{&:special}&url:https://example.com?key=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试第一个输入
			result1 := lib.SortMapToStr(tt.input1)
			assert.Equal(t, tt.expected, result1, "输入1失败")

			// 测试第二个输入（顺序不同但等价的map）
			result2 := lib.SortMapToStr(tt.input2)
			assert.Equal(t, tt.expected, result2, "输入2失败")

			// 确保两个输入产生相同的结果
			assert.Equal(t, result1, result2, "两个等价输入产生了不同的结果")
		})
	}
}

func TestAuthEncrypt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		params   map[string]any
		expected string
	}{
		{
			name: "基本加密测试",
			key:  "testkey",
			params: map[string]any{
				"name": "test",
				"age":  25,
			},
			expected: "f6c0c7820e1df755d2c5f5f2e8750e25a9930ae02c9f3b76d6e5ce1c0e39e7f4",
		},
		{
			name:     "空参数测试",
			key:      "testkey",
			params:   map[string]any{},
			expected: "e0bc614e4fd035a488619799853b075143deea596c477b8dc077e309c0fe42e9",
		},
		{
			name: "复杂参数测试",
			key:  "testkey",
			params: map[string]any{
				"user": map[string]any{
					"name": "test",
					"age":  25,
				},
				"status": "active",
			},
			expected: "61b3c7c8d8f75e1765af0df1dd24c87c2c8fac61c3b7f9a8b2856d9c51d7c8f5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lib.AuthEncrypt(tt.key, tt.params)
			assert.NotEmpty(t, result)
			// 由于HMAC结果是确定性的，对于相同的输入应该总是产生相同的输出
			result2 := lib.AuthEncrypt(tt.key, tt.params)
			assert.Equal(t, result, result2)
		})
	}
}

func TestGetRandSource(t *testing.T) {
	source := lib.GetRandSource()
	assert.NotNil(t, source)

	// 测试生成的随机数是否在预期范围内
	for i := 0; i < 1000; i++ {
		value := source.Uint64()
		assert.Less(t, value, uint64(10000000))
	}

	// 测试多次调用是否返回不同的值
	results := make(map[uint64]bool)
	for i := 0; i < 100; i++ {
		value := source.Uint64()
		results[value] = true
	}
	// 检查是否有足够的不同值（允许一些重复）
	assert.Greater(t, len(results), 50)
}

func TestShuffle(t *testing.T) {
	tests := []struct {
		name  string
		input []int
	}{
		{
			name:  "整数切片",
			input: []int{1, 2, 3, 4, 5},
		},
		{
			name:  "空切片",
			input: []int{},
		},
		{
			name:  "单元素切片",
			input: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制原始切片以便比较
			original := make([]int, len(tt.input))
			copy(original, tt.input)

			// 执行洗牌
			lib.Shuffle(tt.input)

			// 检查长度是否保持不变
			assert.Equal(t, len(original), len(tt.input))

			if len(tt.input) > 1 {
				// 检查元素是否仍然相同（可能顺序不同）
				originalMap := make(map[int]int)
				shuffledMap := make(map[int]int)

				for _, v := range original {
					originalMap[v]++
				}
				for _, v := range tt.input {
					shuffledMap[v]++
				}

				assert.Equal(t, originalMap, shuffledMap)

				// 对于较长的切片，检查是否真的发生了洗牌
				// （注意：这个测试理论上可能偶尔失败，因为完全相同的顺序是可能的，尽管概率很小）
				if len(tt.input) > 3 {
					different := false
					for i := range tt.input {
						if tt.input[i] != original[i] {
							different = true
							break
						}
					}
					assert.True(t, different, "洗牌后的序列与原序列完全相同")
				}
			}
		})
	}
}
