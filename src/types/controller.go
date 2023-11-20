package types

import (
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/text/gstr"
)

// TranslatePlatform 翻译平台
type TranslatePlatform struct {
	md5      string
	Platform string         `json:"platform"`
	Status   int            `json:"status"`
	Level    int            `json:"level"`
	Cfg      map[string]any `json:"cfg"`
	Type     string         `json:"type"`
}

func (t *TranslatePlatform) InitMd5() {
	kArr := garray.NewStrArray()
	for i, _ := range t.Cfg {
		kArr.PushRight(i)
	}
	var s = []string{
		fmt.Sprintf("platform:%s", t.Platform),
		fmt.Sprintf("status:%d", t.Status),
		fmt.Sprintf("level:%d", t.Level),
		fmt.Sprintf("type:%s", t.Type),
	}
	for _, v := range kArr.Sort().Slice() {
		s = append(s, fmt.Sprintf("%s:%s", v, t.Cfg[v]))
	}
	t.md5 = gmd5.MustEncrypt(gstr.Join(s, "-"))
}

func (t *TranslatePlatform) GetMd5() string {
	return t.md5
}
