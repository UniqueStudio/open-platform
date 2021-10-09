package handles

import (
	"context"
	"strconv"

	"github.com/UniqueStudio/open-platform/database"
	"github.com/UniqueStudio/open-platform/pb/sms"
	"github.com/UniqueStudio/open-platform/utils"
	"github.com/xylonx/zapx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type TencentSMSClient struct {
	sms.UnimplementedSMSServiceServer
}

func NewTencentSMSGrpcServer() *TencentSMSClient {
	return &TencentSMSClient{}
}

func (tsc *TencentSMSClient) PushSMS(ctx context.Context, in *sms.PushSMSRequest) (*sms.PushSMSResponse, error) {
	apmCtx, span := utils.Tracer.Start(ctx, "PushSMS")
	defer span.End()

	span.SetAttributes(attribute.Any("requestBody", in))

	signId, err := strconv.ParseUint(in.SignId, 10, 16)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("parse sign id failed", zap.Error(err))
		return nil, err
	}
	templateId, err := strconv.ParseUint(in.TemplateId, 10, 16)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("parse template id failed", zap.Error(err))
		return nil, err
	}

	resp, err := utils.SendMultipleSMS(
		apmCtx, uint(signId), uint(templateId),
		in.PhoneNumber, in.TemplateParamSet,
	)
	if err != nil || resp == nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("push sms failed")
		return nil, err
	}

	data := make([]*sms.SMSStatus, len(*resp))
	for i := range *resp {
		data[i] = &sms.SMSStatus{
			SerialNo:    (*resp)[i].SerialNo,
			PhoneNumber: (*resp)[i].PhoneNumber,
			Fee:         (*resp)[i].Fee,
			ErrCode:     (*resp)[i].ErrCode,
			Message:     (*resp)[i].Message,
		}
	}

	zapx.WithContext(apmCtx).Info("push sms successfully", zap.Any("resp", data))

	return &sms.PushSMSResponse{SMSStatus: data}, nil
}

func (tsc *TencentSMSClient) AddSMSSignature(ctx context.Context, in *sms.AddSMSSignatureRequest) (*sms.UniformResponse, error) {
	apmCtx, span := utils.Tracer.Start(ctx, "AddSMSSignature")
	defer span.End()

	span.SetAttributes(attribute.Any("requestBody", in))

	signs := make([]database.SMSSignature, len(in.Signatures))
	for i := range in.Signatures {
		signs[i].SignID = in.Signatures[i].SignId
		signs[i].Name = in.Signatures[i].SignContent
	}
	resp, err := database.InsertSMSSignatures(&signs)
	if err != nil {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			zapx.WithContext(apmCtx).Error("insert sms signature failed", zap.Error(err))

			return nil, err
		}
	}
	status := make([]*sms.AddStatus, len(*resp))
	for i := range *resp {
		status[i] = &sms.AddStatus{
			Success: true,
			Id:      strconv.FormatUint(uint64((*resp)[i].ID), 10),
			Message: (*resp)[i].Name,
		}
	}
	zapx.WithContext(apmCtx).Info("insert sms signature successfully", zap.Any("status", status))
	return &sms.UniformResponse{Status: status}, nil
}

func (tsc *TencentSMSClient) AddSMSTemplate(ctx context.Context, in *sms.AddSMSTemplateRequest) (*sms.UniformResponse, error) {
	apmCtx, span := utils.Tracer.Start(ctx, "AddSMSTemplate")
	defer span.End()

	span.SetAttributes(attribute.Any("requestBody", in))

	templates := make([]database.SMSTemplate, len(in.Templates))
	for i := range in.Templates {
		templates[i].TemplateID = in.Templates[i].TemplateId
		templates[i].ParamNumber = in.Templates[i].ParamNumber
		templates[i].Content = in.Templates[i].SignContent
	}

	resp, err := database.InsertSMSTemplates(&templates)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("insert template failed", zap.Error(err))

		return nil, err
	}

	status := make([]*sms.AddStatus, len(*resp))
	for i := range *resp {
		status[i] = &sms.AddStatus{
			Success: true,
			Id:      strconv.FormatUint(uint64((*resp)[i].ID), 10),
			Message: (*resp)[i].Content,
		}
	}
	zapx.WithContext(apmCtx).Info("insert sms signature successfully", zap.Any("status", status))
	return &sms.UniformResponse{Status: status}, nil
}
