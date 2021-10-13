package handles

import (
	"github.com/UniqueStudio/open-platform/database"
	"github.com/UniqueStudio/open-platform/pkg"
	"github.com/UniqueStudio/open-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"net/http"
)

func GetEmailTemplateHandler(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "AddSMSTemplate")
	defer span.End()

	templates, err := database.GetAllEmailTemplate(ctx)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get all templates failed")

		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("templates", templates))
	zapx.WithContext(apmCtx).Info("get all email templates sufficiently")
	ctx.JSON(http.StatusOK, pkg.SuccessResponse(templates))
}

func SendSingleEmailHandler(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "SendSingleEmailTemplate")
	defer span.End()

	data := new(pkg.SingleEmailReq)
	if err := ctx.ShouldBindJSON(data); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("requestBody", data))
	zapx.WithContext(apmCtx).Info("bind JSON sufficiently")

	resp, err := utils.SendSingleEmail(apmCtx, data.EmailTo, data.TemplateID, data.Subject, data.TemplateParamSet)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send single email failed")
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}
	span.SetAttributes(attribute.Any("send email Response", resp))
	zapx.WithContext(apmCtx).Info("send single email sufficiently")
	ctx.JSON(http.StatusOK, pkg.SuccessResponse(resp))
}

func SendGroupEmailHandler(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "SendGroupEmail")
	defer span.End()

	data := new(pkg.MultipleEmailReq)
	if err := ctx.ShouldBindJSON(data); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("requestBody", data))
	zapx.WithContext(apmCtx).Info("bind JSON sufficiently")
	resp, err := utils.SendMultipleEmail(apmCtx, data.EmailTo, data.TemplateID, data.Subject, data.TemplateParamSet)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send single email failed")
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("send group email Response", resp))
	zapx.WithContext(apmCtx).Info("send group email sufficiently")
	ctx.JSON(http.StatusOK, pkg.SuccessResponse(resp))
}

func InsertEmailTemplateHandler(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "AddEmailTemplate")
	defer span.End()

	data := pkg.EmailTemplatesReq{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("requestBody", data))
	zapx.WithContext(apmCtx).Info("bind JSON sufficiently")

	template := database.EmailTemplate{}
	template.TemplateID = data.TemplateID
	template.Content = data.Content
	template.ParamNumber = data.ParamNumber

	resp, err := database.InsertEmailTemplate(ctx, &[]database.EmailTemplate{template})

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send multiple emails failed")
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("sending email Response", resp))
	zapx.WithContext(apmCtx).Info("send multiple email sufficiently")
	ctx.JSON(http.StatusOK, pkg.SuccessResponse(resp))
}
