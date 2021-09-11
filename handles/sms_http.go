// Deprecated: the http handler will be removed soon.
package handles

import (
	"net/http"

	"github.com/UniqueStudio/open-platform/config"
	"github.com/UniqueStudio/open-platform/database"
	"github.com/UniqueStudio/open-platform/pkg"
	"github.com/UniqueStudio/open-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func AddSMSSign(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "AddSMSSign")
	defer span.End()

	data := []pkg.AddSignReq{}
	if err := ctx.ShouldBindJSON(&data); err != nil || data == nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("requestBody", data))
	zapx.WithContext(apmCtx).Info("bind JSON sufficiently")

	signs := make([]database.SMSSignature, len(data))
	for i := range data {
		signs[i].SignID = data[i].SignID
		signs[i].Name = data[i].Name
	}

	resp, err := database.InsertSMSSignatures(&signs)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}

	zapx.WithContext(apmCtx).Info("insert sms signs successfully", zap.Any("signatures", resp))

	ctx.JSON(http.StatusOK, pkg.SuccessResponse(resp))
}

func AddSMSTemplate(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "AddSMSTemplate")
	defer span.End()

	data := []pkg.AddTemplateReq{}
	if err := ctx.ShouldBindJSON(&data); err != nil || data == nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("requestBody", data))
	zapx.WithContext(apmCtx).Info("bind JSON sufficiently")

	templates := make([]database.SMSTemplate, len(data))
	for i := range data {
		templates[i].TemplateID = data[i].TemplateID
		templates[i].ParamNumber = data[i].ParamNumber
		templates[i].Content = data[i].Content
	}

	resp, err := database.InsertSMSTemplates(&templates)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}

	zapx.WithContext(apmCtx).Info("insert templates successfully", zap.Any("templates", resp))

	ctx.JSON(http.StatusOK, pkg.SuccessResponse(resp))
}

func SendSingleSMSHandler(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "SingleSMSHandler")
	defer span.End()

	data := new(pkg.SingleSMSReq)
	if err := ctx.ShouldBindJSON(data); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("requestBody", data))
	zapx.WithContext(apmCtx).Info("bind JSON sufficiently")

	signId := config.Config.Tencent.SMS.DefaultVirtualSignId
	if data.SignID != nil {
		signId = *data.SignID
	}
	resp, err := utils.SendSingleSMS(
		apmCtx, signId, data.TemplateID,
		data.PhoneNumber, data.TemplateParamSet,
	)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send single sms failed")

		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("tencentResponse", resp))
	zapx.WithContext(apmCtx).Info("send single sms sufficiently")
	ctx.JSON(http.StatusOK, pkg.SuccessResponse(resp))
}

func SendGroupSMSHandler(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "MultipleSMSHandler")
	defer span.End()

	data := new(pkg.MultipleSMSReq)
	if err := ctx.ShouldBindJSON(data); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind JSON failed", zap.Error(err))

		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("requestBody", data))
	zapx.WithContext(apmCtx).Info("bind JSON sufficiently")

	signId := config.Config.Tencent.SMS.DefaultVirtualSignId
	if data.SignID != nil {
		signId = *data.SignID
	}
	resp, err := utils.SendMultipleSMS(
		apmCtx, signId, data.TemplateID,
		data.PhoneNumber, data.TemplateParamSet,
	)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send multiple sms failed")

		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("tencentResponse", resp))
	zapx.WithContext(apmCtx).Info("send single sms sufficiently")
	ctx.JSON(http.StatusOK, pkg.SuccessResponse(resp))
}

func GetSMSTemplateHandler(ctx *gin.Context) {
	apmCtx, span := utils.Tracer.Start(ctx.Request.Context(), "GetAllSMSTemplates")
	defer span.End()

	templates, err := database.GetAllTemplate(apmCtx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("get all templates failed")

		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err))
		return
	}

	span.SetAttributes(attribute.Any("templates", templates))
	zapx.WithContext(apmCtx).Info("get all templates sufficiently")
	ctx.JSON(http.StatusOK, pkg.SuccessResponse(templates))
}
