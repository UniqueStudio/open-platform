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
	initSMSRouter(smsrouter)
}

func initSMSRouter(r *gin.RouterGroup) {
	// need auth
	r.Use(middleware.Auth())
	r.POST("/send_single", handles.SendSingleSMSHandler)
	r.POST("/send_group", handles.SendGroupSMSHandler)
	r.GET("/templates", handles.GetSMSTemplateHandler)
}
