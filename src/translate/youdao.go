package translate

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/xgd16/gf-x-tool/xlib"
)

// YouDaoTranslate 有道翻译
func YouDaoTranslate(YouDaoConfig *YouDaoConfigType, from, to, text string) (result []string, fromLang string, err error) {
	if YouDaoConfig == nil || YouDaoConfig.AppKey == "" || YouDaoConfig.Url == "" || YouDaoConfig.SecKey == "" {
		err = errors.New("有道翻译配置异常")
		return
	}
	// 语言标记转换
	from, err = SafeLangType(from, YouDaoTranslateMode)
	if err != nil {
		return
	}
	if to == "auto" {
		err = errors.New("转换后语言不能为auto")
		return
	}

	truncate := func(s string) string {
		l := gstr.LenRune(s)
		if l <= 20 {
			return s
		}
		return fmt.Sprintf("%s%d%s", gstr.SubStrRune(s, 0, 10), l, gstr.SubStrRune(s, l-10, l))
	}

	salt := gtime.Now().UnixMilli()
	curTime := int(math.Round(float64(salt / 1000)))
	signStr := fmt.Sprintf("%s%s%d%d%s", YouDaoConfig.AppKey, truncate(text), salt, curTime, YouDaoConfig.SecKey)
	sign := xlib.Sha256(signStr)

	post, err := g.Client().SetTimeout(time.Duration(YouDaoConfig.CurlTimeOut)*time.Millisecond).Post(gctx.New(), YouDaoConfig.Url, g.Map{
		"q":        text,
		"appKey":   YouDaoConfig.AppKey,
		"salt":     salt,
		"from":     from,
		"to":       to,
		"sign":     sign,
		"signType": "v3",
		"curtime":  curTime,
	})
	if err != nil {
		return
	}
	defer func() { _ = post.Close() }()
	postResp := post.ReadAllString()
	if post.StatusCode != 200 {
		err = fmt.Errorf("请求失败 状态码: %d 返回结果: %s", post.StatusCode, postResp)
		return
	}
	json, err := gjson.DecodeToJson(postResp)
	if err != nil {
		return
	}
	if json.Get("errorCode").Int() != 0 {
		err = fmt.Errorf("请求失败errorCode: %d err: %s", json.Get("errorCode").Int(), postResp)
		return
	}
	// 获取 from
	returnFrom := gstr.Split(json.Get("l").String(), "2")[0]
	return json.Get("translation").Strings(), returnFrom, nil
}
