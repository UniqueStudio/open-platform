package main

import (
	_ "./docs"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"open-platform/handler"
	"open-platform/middleware"
	"open-platform/utils"
)

// @title Open Platform API
// @version 0.1
// @description This is a Unique Studio Open Platform API server.

// @contact.name Fred Liang
// @contact.url https://blog.fredliang.cn
// @contact.email info@fredliang.cn

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Token

// @license.name MPL-2.0

// @host open.hustunique.com
// @schemes https
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
		token.GET("/template", handler.GetSMSTemplate)
		token.GET("/department", handler.GetDepartmentListHandler)
		token.GET("/department/:departmentID", handler.GetDepartmentUsersHandler)
	}

	r.GET("/department/:departmentID", handler.GetDepartmentUsersHandler)
	r.GET("/department", handler.GetDepartmentListHandler)

	// Use ginSwagger gen api doc
	r.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusPermanentRedirect, "/doc/index.html") })

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
