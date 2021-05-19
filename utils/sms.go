package utils

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
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
	request.PhoneNumberSet = common.StringPtrs([]string{phoneNumber})
	request.TemplateParamSet = common.StringPtrs(templateParamSet)
	request.TemplateID = common.StringPtr(templateID)
	request.SmsSdkAppid = common.StringPtr(AppConfig.TencentCloudSDKSMS.SDKAppID)
	//request.Sign = common.StringPtr(sign)

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

