package router

import (
	"net/http"

	"github.com/UniqueStudio/open-platform/handles"
	"github.com/UniqueStudio/open-platform/middleware"
	"github.com/UniqueStudio/open-platform/pkg"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	middleware.InitGlobalMiddleWare(r)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, pkg.SuccessResponse("pong"))
	})

	smsrouter := r.Group("/sms")
	emailrouter := r.Group("/email")

	initSMSRouter(smsrouter)
	initEmailRouter(emailrouter)
}

func initSMSRouter(r *gin.RouterGroup) {
	// need auth
	r.Use(middleware.Authentication())
	r.Use(middleware.Authorization())
	// get sms infos
	r.POST("/send_single", handles.SendSingleSMSHandler)
	r.POST("/send_group", handles.SendGroupSMSHandler)
	r.GET("/templates", handles.GetSMSTemplateHandler)
}

func initEmailRouter(r *gin.RouterGroup) {
	r.Use(middleware.Authentication())
	r.Use(middleware.Authorization())

	r.GET("/templates", handles.GetEmailTemplateHandler)
	r.POST("/send_single", handles.SendSingleEmailHandler)
	r.POST("/send_group", handles.SendGroupEmailHandler)
	r.POST("/templates", handles.InsertEmailTemplateHandler)
}
