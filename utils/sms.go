package utils

import (
	"fmt"
	"time"

	qcloudsms "github.com/fredliang44/qcloudsms_go"
)

// SendQCSMS is a func to handle sms with qcloud
func SendQCSMS(Phone string, Template int, ParamList []string) (isOK bool, message string, errID string) {
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

// SendQCSMSMulti is a func to handle sms with qcloud
func SendQCSMSMulti(PhoneList []string, Template int, ParamList []string) (isOK bool, message string, errID string) {
	opt := qcloudsms.NewOptions(AppConfig.QcloudSMS.AppID, AppConfig.QcloudSMS.AppKey, AppConfig.QcloudSMS.Sign)

	var client = qcloudsms.NewClient(opt)
	client.SetDebug(true)
	TelList := []qcloudsms.SMSTel{}
	for _, phone := range PhoneList {
		TelList = append(TelList, qcloudsms.SMSTel{Nationcode: "86", Mobile: phone})
	}

	var t = qcloudsms.SMSMultiReq{
		Params: ParamList,
		Tel:    TelList,
		Sign:   AppConfig.QcloudSMS.Sign,
		TplID:  uint(Template),
	}

	isOK, err := client.SendSMSMulti(t)
	return isOK, fmt.Sprintln(err), fmt.Sprintln(err)
}

// GetQCSMSTemplate is a func to Get Qcloud SMS Template
func GetQCSMSTemplate() qcloudsms.TemplateGetResult {
	opt := qcloudsms.NewOptions(AppConfig.QcloudSMS.AppID, AppConfig.QcloudSMS.AppKey, AppConfig.QcloudSMS.Sign)

	opt.Debug = true

	var client = qcloudsms.NewClient(opt)

	Template, _ := client.GetTemplateByPage(0, 30)
	return Template
}

// AddQCSMSTemplate is a func to add timepla
func AddQCSMSTemplate(title, text, remark string) (qcloudsms.TemplateResult, error) {
	opt := qcloudsms.NewOptions(AppConfig.QcloudSMS.AppID, AppConfig.QcloudSMS.AppKey, AppConfig.QcloudSMS.Sign)

	opt.Debug = true

	var client = qcloudsms.NewClient(opt)

	TemplateResult, err := client.NewTemplate(qcloudsms.TemplateNew{
		Title:  title,
		Text:   text,
		Type:   0,
		Remark: remark,
		Time:   time.Now().Unix(),
	})

	return TemplateResult, err
}

// GetQCSMSTemplateStatus is a func to get Get QC SMS Template Status
func GetQCSMSTemplateStatus(id []uint) (qcloudsms.TemplateGetResult, error) {
	opt := qcloudsms.NewOptions(AppConfig.QcloudSMS.AppID, AppConfig.QcloudSMS.AppKey, AppConfig.QcloudSMS.Sign)

	opt.Debug = true

	var client = qcloudsms.NewClient(opt)

	templateGetResult, err := client.GetTemplateByID(id)

	return templateGetResult, err
}
