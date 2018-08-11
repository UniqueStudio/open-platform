package main

import (
	"fmt"
	"open-platform/handler"
	"open-platform/middleware"
	"open-platform/utils"

	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())

	// CORS support
	r.Use(middleware.CORSMiddleware())

	// Recovery from internal server error
	r.Use(nice.Recovery(handler.RecoveryHandler))

	store := cookie.NewStore([]byte(utils.AppConfig.Server.SecretKey))
	r.Use(sessions.Sessions("Status", store))

	r.LoadHTMLGlob("static/*/*.tmpl")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/login", handler.LoginHandler)
	r.GET("/login/:app", handler.LoginHandler)
	r.GET("/check", handler.CheckAuthorityHandler)

	r.GET("/api", handler.GenAccessKeyHandler)

	app := r.Group("/app")
	app.Use(middleware.Login())
	{
		app.GET("/:app", handler.RenderAppStaticFilesHandler)
		// TODO app.StaticFS("/message/", http.Dir("./static/message/dist/"))
	}

	login := r.Group("/")
	login.Use(middleware.Login())
	{
		login.GET("/test", handler.TestHandler)
		login.GET("/decode", handler.DecodeHandler)
	}

	r.GET("/auth", handler.AuthHandler)
	r.GET("/auth/:app", handler.AuthAPPHandler)

	weixin := r.Group("/weixin")
	weixin.Use(middleware.Auth())
	{
		weixin.POST("/sms", handler.SendSMSHandler)
		weixin.GET("/sms/template", handler.GetSMSTemplateHandler)
		weixin.POST("/sms/template", handler.AddSMSTemplateHandler)
		weixin.GET("/sms/template/:id", handler.GetSMSTemplateStatusHandler)
		weixin.GET("/department", handler.GetDepartmentListHandler)
		weixin.GET("/department/:departmentID", handler.GetDepartmentUsersHandler)
	}

	message := r.Group("/message")
	message.Use(middleware.Auth())
	{
		message.POST("/sms", handler.SendSMSHandler)
		message.GET("/sms/template", handler.GetSMSTemplateHandler)
		message.POST("/sms/template", handler.AddSMSTemplateHandler)
		message.GET("/sms/template/:id", handler.GetSMSTemplateStatusHandler)
		message.POST("/mail", handler.SendMailHandler)
	}

	showStatus()

	// Run Server
	r.Run(utils.AppConfig.Server.Host + ":" + utils.AppConfig.Server.Port)
}

func showStatus() {
	fmt.Println("\n===================================" +
		"\nAPP         : " + utils.AppConfig.APPName +
		"\nRunning On  : " + utils.AppConfig.Server.Host + ":" + utils.AppConfig.Server.Port +
		"\n===================================")
}