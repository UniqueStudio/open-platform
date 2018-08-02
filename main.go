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

	token := r.Group("/weixin")
	token.Use(middleware.Auth())
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
