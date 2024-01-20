package translate

import (
	baseErr "errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"github.com/xgd16/gf-x-tool/xtranslate"
)

func TencentTranslate(config *TencentConfigType, from, to, text string) (result []string, fromLang string, err error) {
	if config.SecretId == "" || config.SecretKey == "" || config.Region == "" || config.Url == "" {
		err = baseErr.New("腾讯翻译配置异常")
		return
	}

	from, err = xtranslate.SafeLangType(from, TencentTranslateMode)
	to, err = xtranslate.SafeLangType(to, TencentTranslateMode)
	if err != nil {
		return
	}

	credential := common.NewCredential(
		config.SecretId,
		config.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = config.Url

	client, err := tmt.NewClient(credential, config.Region, cpf)
	if err != nil {
		return
	}
	request := tmt.NewTextTranslateRequest()

	request.SourceText = common.StringPtr(text)
	request.Source = common.StringPtr("auto")
	request.Target = common.StringPtr(to)
	request.ProjectId = common.Int64Ptr(0)

	response, err := client.TextTranslate(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		err = baseErr.New(fmt.Sprintf("An API error has returned: %s", err))
		return
	}
	if err != nil {
		return
	}
	jsonData, err := gjson.DecodeToJson(response.ToJsonString())
	if err != nil {
		return
	}
	result = jsonData.Get("Response.TargetText").Strings()
	fromLang, err = xtranslate.GetYouDaoLang(jsonData.Get("Response.Source").String(), TencentTranslateMode)
	return
}
