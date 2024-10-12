package lib

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"sort"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

func AuthEncrypt(key string, params map[string]any) string {
	return gmd5.MustEncryptString(key + SortMapToStr(params))
}

// SortMapToStr 排序 map 结构数据使其有序化
func SortMapToStr(data map[string]any) (str string) {
	var mapArr = make([]string, 0)
	for k, v := range data {
		if vItem, ok := v.(map[string]any); ok {
			mapArr = append(mapArr, fmt.Sprintf("%s:|%s|", k, SortMapToStr(vItem)))
			continue
		}
		valueType := reflect.TypeOf(v).Kind()
		if valueType == reflect.Array || valueType == reflect.Slice {
			mapArr = append(mapArr, fmt.Sprintf("%s:%s", k, gstr.JoinAny(v, ",")))
			continue
		}
		mapArr = append(mapArr, fmt.Sprintf("%s:%s", k, gconv.String(v)))
	}
	sort.Strings(mapArr)
	str = gstr.Join(mapArr, "&")
	return
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
	n := gconv.Uint64(grand.N(1, 9999999))
	// g.DumpWithType(n)
	return n
}

func GetRandSource() *RandSource {
	return &RandSource{}
}

func Shuffle[T any](slice []T) {
	r := rand.New(GetRandSource())
	r.Shuffle(len(slice), func(i, j int) {
		// g.DumpWithType(i, j)
		slice[i], slice[j] = slice[j], slice[i]
	})
}
