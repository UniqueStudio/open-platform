package utils

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func SendSingleSms(phoneNumber string, templateParamSet []string, templateID string) (*sms.SendSmsResponse, error) {
	credential := common.NewCredential(
		AppConfig.TencentCloudSDKSMS.SecretID,
		AppConfig.TencentCloudSDKSMS.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "", cpf)

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(AppConfig.TencentCloudSDKSMS.SDKAppID)
	request.SignName = common.StringPtr(AppConfig.TencentCloudSDKSMS.Sign.Content)
	request.TemplateParamSet = common.StringPtrs(templateParamSet)
	request.TemplateId = common.StringPtr(templateID)
	request.PhoneNumberSet = common.StringPtrs([]string{phoneNumber})

	response, err := client.SendSms(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	//fmt.Printf("send success, res: %s\n", response.ToJsonString())
	return response, nil
}

func GetTemplates() []*SMSTemplate {
	return AppConfig.TencentCloudSDKSMS.Templates
}
