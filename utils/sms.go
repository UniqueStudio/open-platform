package utils

import (
	"fmt"

	qcloudsms "github.com/qichengzx/qcloudsms_go"
)

// SendQCSMS is a func to handle sms with qcloud
func SendQCSMS(Phone string, Template int, ParamList []string) (isOK bool, msg string, errID string) {
	opt := qcloudsms.NewOptions(AppConfig.QcloudSMS.AppID, AppConfig.QcloudSMS.AppKey, AppConfig.QcloudSMS.Sign)

	var client = qcloudsms.NewClient(opt)
	client.SetDebug(true)

	var t = qcloudsms.SMSSingleReq{
		Params: ParamList,
		Tel:    qcloudsms.SMSTel{Nationcode: "86", Mobile: Phone},
		Sign:   AppConfig.QcloudSMS.Sign,
		TplID:  Template,
	}

	isOK, err := client.SendSMSSingle(t)
	return isOK, fmt.Sprintln(err), fmt.Sprintln(err)
}

func GetQCSMSTemplate() qcloudsms.TemplateGetResult {
	opt := qcloudsms.NewOptions(AppConfig.QcloudSMS.AppID, AppConfig.QcloudSMS.AppKey, AppConfig.QcloudSMS.Sign)

	opt.Debug = true

	var client = qcloudsms.NewClient(opt)

	Template, _ :=client.GetTemplateByPage(0, 30)
	return Template
}
