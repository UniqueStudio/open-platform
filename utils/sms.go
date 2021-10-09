package utils

import (
	"context"
	"strings"

	"github.com/UniqueStudio/open-platform/config"
	"github.com/UniqueStudio/open-platform/database"
	"github.com/UniqueStudio/open-platform/pkg"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/xylonx/zapx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

var (
	SMSClient *sms.Client
)

// the client is static struct.
// Therefore it is good to cache it in memory instead of create it each time
func SetupSMSClient() error {
	credential := common.NewCredential(
		config.Config.Tencent.SMS.SecretID,
		config.Config.Tencent.SMS.SecretKey,
	)

	cpf := profile.NewClientProfile()
	client, err := sms.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return err
	}
	SMSClient = client
	return nil
}

func SendSingleSMS(ctx context.Context, vsignId, vtemplateId uint, phone string, paramSet []string) (*[]pkg.SMSSendStatus, error) {
	apmCtx, span := Tracer.Start(ctx, "SendSingleSMS")
	defer span.End()
	phones := addPhoneAreaCodePrefix([]string{phone})
	return sendTencentSMS(apmCtx, vsignId, vtemplateId, paramSet, phones)
}

func SendMultipleSMS(ctx context.Context, vsignId, vtemplateId uint, phones, paramSet []string) (*[]pkg.SMSSendStatus, error) {
	apmCtx, span := Tracer.Start(ctx, "SendSingleSMS")
	defer span.End()
	phones = addPhoneAreaCodePrefix(phones)
	return sendTencentSMS(apmCtx, vsignId, vtemplateId, paramSet, phones)
}

func sendTencentSMS(ctx context.Context, vsignId, vtemplateId uint, paramSet, phones []string) (*[]pkg.SMSSendStatus, error) {
	apmCtx, span := Tracer.Start(ctx, "SendTencentSMS")
	defer span.End()

	signature, err := database.GetSMSSignatureByVirtualID(apmCtx, vsignId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		zapx.WithContext(apmCtx).Error("get sms signature from virtualId failed", zap.Error(err))
		return nil, err
	}
	span.SetAttributes(attribute.Any("SMSSignature", signature))

	template, err := database.GetSMSTemplateByVirtualID(apmCtx, vtemplateId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		zapx.WithContext(apmCtx).Error("get sms template from virtualId failed", zap.Error(err))
		return nil, err
	}
	span.SetAttributes(attribute.Any("SMSTemplate", template))

	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = common.StringPtr(config.Config.Tencent.SMS.SDKAppID)
	req.SignName = common.StringPtr(signature.Name)
	req.TemplateParamSet = common.StringPtrs(paramSet)
	req.TemplateId = common.StringPtr(template.TemplateID)
	req.PhoneNumberSet = common.StringPtrs(phones)

	span.SetAttributes(attribute.Any("SMSRequest", req))

	resp, err := SMSClient.SendSms(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		zapx.WithContext(apmCtx).Error("send sms request error")
		return nil, err
	}

	return pkg.TencentSMSToSMSResp(resp), nil
}

// TODO:
func GetTemplates() *[]database.SMSTemplate {
	return nil
}

// add area code in front of phone number.
// default area code is +86
// TODO: pass pointer to optimize
func addPhoneAreaCodePrefix(phones []string) []string {
	nphones := make([]string, len(phones))
	for i := range phones {
		if strings.HasPrefix(phones[i], "+") {
			nphones[i] = phones[i]
		} else {
			nphones[i] = "+86" + phones[i]
		}
	}
	return nphones
}
