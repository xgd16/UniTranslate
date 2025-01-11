package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func AuthEncrypt(key string, params map[string]any) string {
	data := SortMapToStr(params)
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SortMapToStr 将map转换为有序字符串，确保深层数据一致性
func SortMapToStr(data map[string]any) string {
	if len(data) == 0 {
		return ""
	}
	
	// 预分配容量以提高性能
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.Grow(len(data) * 32) // 预估字符串长度

	for i, k := range keys {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(k)
		b.WriteByte(':')
		
		// 处理不同类型的值
		switch v := data[k].(type) {
		case map[string]any:
			// 递归处理嵌套的map
			b.WriteByte('{')
			b.WriteString(SortMapToStr(v))
			b.WriteByte('}')
		case []any:
			// 处理数组/切片
			b.WriteByte('[')
			for j, item := range v {
				if j > 0 {
					b.WriteByte(',')
				}
				switch iv := item.(type) {
				case map[string]any:
					b.WriteByte('{')
					b.WriteString(SortMapToStr(iv))
					b.WriteByte('}')
				default:
					b.WriteString(gconv.String(iv))
				}
			}
			b.WriteByte(']')
		default:
			b.WriteString(gconv.String(v))
		}
	}
	
	return b.String()
}

func SqliteTableIsExists(db gdb.DB, tableName string) (isExists bool, err error) {
	count, err := db.Model("sqlite_master").Count(g.Map{
		"type": "table",
		"name": tableName,
	})
	if err != nil {
		return
	}
	isExists = count > 0
	return
}

type RandSource struct {
}

func (t *RandSource) Uint64() uint64 {
	n := gconv.Uint64(rand.Int63n(9999999))
	// g.DumpWithType(n)
	return n
}

func GetRandSource() *RandSource {
	return &RandSource{}
}

func Shuffle[T any](slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
