package translate

import (
	baseErr "errors"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
)

type TencentConfigType struct {
	Url       string `json:"url"`
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
	Region    string `json:"region"`
}

func (t *TencentConfigType) Translate(req *TranslateReq) (resp []*TranslateResp, err error) {
	if t == nil || t.SecretId == "" || t.SecretKey == "" || t.Region == "" || t.Url == "" {
		err = baseErr.New("腾讯翻译配置异常")
		return
	}
	mode := t.GetMode()
	from, err := SafeLangType(req.From, mode)
	if err != nil {
		return
	}
	to, err := SafeLangType(req.To, mode)
	if err != nil {
		return
	}

	credential := common.NewCredential(
		t.SecretId,
		t.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = t.Url

	client, err := tmt.NewClient(credential, t.Region, cpf)
	if err != nil {
		return
	}
	request := tmt.NewTextTranslateBatchRequest()

	request.SourceTextList = common.StringPtrs(req.Text)
	request.Source = common.StringPtr(from)
	request.Target = common.StringPtr(to)
	request.ProjectId = common.Int64Ptr(0)

	response, err := client.TextTranslateBatch(request)
	if err != nil {
		return
	}
	var tencentCloudSDKError *errors.TencentCloudSDKError
	if baseErr.As(err, &tencentCloudSDKError) {
		err = fmt.Errorf("an API error has returned: %s", tencentCloudSDKError.Error())
		return
	}
	lang, err := GetYouDaoLang(*response.Response.Source, mode)
	if err != nil {
		return
	}
	for _, item := range response.Response.TargetTextList {
		resp = append(resp, &TranslateResp{
			Text:     *item,
			FromLang: lang,
		})
	}
	return
}

func (t *TencentConfigType) GetMode() string {
	return TencentTranslateMode
}
