package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"open-platform/handler"
	"open-platform/middleware"
	"open-platform/utils"
)

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte(utils.AppConfig.Server.SecretKey))
	r.Use(sessions.Sessions("Status", store))

	r.LoadHTMLGlob("static/html/*")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/login", handler.LoginHandler)
	r.GET("/check", handler.CheckAuthorityHandler)

	r.GET("/api", handler.GenAccessKeyHandler)

	login := r.Group("/")
	login.Use(middleware.Login())
	{
		login.GET("/test", handler.TestHandler)
		login.GET("/decode", handler.DecodeHandler)
	}

	r.GET("/auth", handler.AuthHandler)

	token := r.Group("/weixin")
	token.Use(middleware.Admin())
	{
		token.POST("/sms", handler.SendSMSHandler)
		token.GET("/sms/template", handler.GetSMSTemplateHandler)
		token.GET("/department", handler.GetDepartmentListHandler)
		token.GET("/department/:departmentID", handler.GetDepartmentUsersHandler)
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
