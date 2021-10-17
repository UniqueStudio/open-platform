package handles

import (
	"context"
	"github.com/UniqueStudio/open-platform/database"
	"github.com/UniqueStudio/open-platform/pb/email"
	"github.com/UniqueStudio/open-platform/utils"
	"github.com/xylonx/zapx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type LarkEmailClient struct {
	email.UnimplementedEmailServiceServer
}

func NewLarkEmailClient() *LarkEmailClient {
	return &LarkEmailClient{}
}

func (lec *LarkEmailClient) PushEmail(ctx context.Context, req *email.PushEmailRequest) (*email.PushEmailResponse, error) {
	apmCtx, span := utils.Tracer.Start(ctx, "PushEmail")
	defer span.End()

	span.SetAttributes(attribute.Any("requestBody", req))
	set := make([][]string, len(req.To))
	template := database.GetEmailTemplateByID(ctx, uint(req.GetTemplateID()))
	index := 0
	for i := 0; i < len(req.TemplateParamSet); i += int(template.ParamNumber) {
		set[index] = req.TemplateParamSet[i : i+int(template.ParamNumber)]
		index++
	}

	resp, err := utils.SendMultipleEmail(apmCtx, req.To, uint(req.TemplateID), req.Subject, set)

	if err != nil || resp == nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send multiple email fail", zap.Error(err))
		return nil, err
	}

	data := make([]*email.EmailStatus, 0)

	for _, v := range resp {
		data = append(data, &email.EmailStatus{
			To:     v.To,
			Err:    v.Err != "",
			ErrMsg: v.Err,
		})
	}

	zapx.WithContext(apmCtx).Info("push email successfully", zap.Any("resp", data))

	return &email.PushEmailResponse{EmailStatus: data},nil
}

func (lec *LarkEmailClient) AddEmailTemplate(ctx context.Context, in *email.AddEmailTemplateRequest) (*email.AddEmailTemplateResponse, error) {
	apmCtx, span := utils.Tracer.Start(ctx, "AddEmailTemplate")
	defer span.End()

	span.SetAttributes(attribute.Any("requestBody", in))

	templates := make([]database.EmailTemplate, len(in.EmailTemplate))
	for i := range in.EmailTemplate {
		templates[i].TemplateID = in.EmailTemplate[i].GetTemplateID()
		templates[i].ParamNumber = in.EmailTemplate[i].GetParamNumber()
		templates[i].Content = in.EmailTemplate[i].GetContent()
		templates[i].TemplateName = in.EmailTemplate[i].GetTemplateName()
	}

	_, err := database.InsertEmailTemplate(apmCtx, &templates)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("insert template failed", zap.Error(err))

		return &email.AddEmailTemplateResponse{}, err
	}

	_template := make([]*email.EmailTemplate, 0)
	for _, v := range templates {
		_template = append(_template, &email.EmailTemplate{
			Content:      v.Content,
			ParamNumber:  v.ParamNumber,
			TemplateID:   v.TemplateID,
			TemplateName: v.TemplateName,
		})
	}

	return &email.AddEmailTemplateResponse{Success: true, Templates: _template}, nil

}

func (lec *LarkEmailClient) GetAllEmailTemplate(ctx context.Context, req *email.GetAllTemplatesRequest) (*email.GetAllTemplatesResponse, error) {
	apmCtx, span := utils.Tracer.Start(ctx, "GetAllEmailTemplate")
	defer span.End()

	span.SetAttributes(attribute.Any("requestBody", req))

	templates, err := database.GetAllEmailTemplate(ctx)
	_templates := make([]*email.EmailTemplate, 0)
	for _, v := range *templates {
		_templates = append(_templates, &email.EmailTemplate{
			Content:      v.Content,
			ParamNumber:  v.ParamNumber,
			TemplateID:   v.TemplateID,
			TemplateName: v.TemplateName,
		})
	}
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("insert template failed", zap.Error(err))

		return &email.GetAllTemplatesResponse{EmailTemplate: _templates, Err: true, ErrMsg: err.Error()}, err
	}
	return &email.GetAllTemplatesResponse{EmailTemplate: _templates}, nil

}
