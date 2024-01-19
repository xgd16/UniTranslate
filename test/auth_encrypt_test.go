package test

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"uniTranslate/src/lib"
)

var paramsData = g.Map{
	"c": g.Map{
		"cc": 1,
		"cb": 2,
		"ca": 3,
		"cd": 4,
	},
	"a": 1,
	"b": []int{4, 1, 2},
}

func TestSortMapToStr(t *testing.T) {
	gtest.Assert("a:1&b:4,1,2&c:|ca:3&cb:2&cc:1&cd:4|", lib.SortMapToStr(paramsData))
}

func TestAuthEncrypt(t *testing.T) {
	gtest.Assert(lib.AuthEncrypt("123456", paramsData), "1ccbf6fc835046b4716b38f5dbd9c75a")
}
