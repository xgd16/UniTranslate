package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

// BaiduConfigType 百度的配置类型
type BaiduConfigType struct {
	CurlTimeOut int    `json:"curlTimeOut"`
	Url         string `json:"url"`
	AppId       string `json:"appId"`
	Key         string `json:"key"`
}

type BaiduHTTPTranslateResp struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []TransResult `json:"trans_result"`
}

type TransResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func (t *BaiduConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	if t == nil || t.Url == "" || t.AppId == "" || t.Key == "" {
		err = errors.New("百度翻译配置异常")
		return
	}
	mode := t.GetMode()
	// 语言标记转换
	from, err := SafeLangType(req.From, mode)
	if err != nil {
		return
	}
	to, err := SafeLangType(req.To, mode)
	if err != nil {
		return
	}
	salt := gtime.Now().UnixMilli()
	signStr := fmt.Sprintf("%s%s%d%s", t.AppId, req.TextStr, salt, t.Key)
	sign, err := gmd5.EncryptString(signStr)
	// 处理MD5加密失败
	if err != nil {
		return
	}
	// 发起请求
	post, err := g.Client().SetTimeout(time.Duration(t.CurlTimeOut)*time.Millisecond).Post(gctx.New(), t.Url, g.Map{
		"q":     req.TextStr,
		"from":  from,
		"to":    to,
		"appid": t.AppId,
		"salt":  salt,
		"sign":  sign,
	})
	// 处理请求失败
	if err != nil {
		return
	}
	// 推出函数时关闭链接
	defer func() { _ = post.Close() }()
	// 返回的json解析
	respByte := post.ReadAll()
	// 判断状态码
	if post.StatusCode != 200 {
		err = fmt.Errorf("请求失败 状态码: %d 返回结果: %s", post.StatusCode, respByte)
		return
	}
	httpResp := new(BaiduHTTPTranslateResp)
	if err = json.Unmarshal(respByte, httpResp); err != nil {
		return
	}
	resp = make([]*TranslateResp, 0)
	for _, item := range httpResp.TransResult {
		lang, err1 := GetYouDaoLang(httpResp.From, mode)
		if err1 != nil {
			err = err1
			return
		}
		resp = append(resp, &TranslateResp{
			Text:     item.Dst,
			FromLang: lang,
		})
	}

	return
}

func (t *BaiduConfigType) GetMode() string {
	return BaiduTranslateMode
}
