package main


import (
	"fmt"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"open-platform/handler"
	"open-platform/middleware"
	"open-platform/utils"
)

func main() {

	r := gin.New()
	r.Use(gin.Logger())

	// CORS support
	r.Use(middleware.CORSMiddleware())

	// Recovery from internal server error
	r.Use(nice.Recovery(handler.RecoveryHandler))

	// Static files
	r.Use(static.Serve("/static", static.LocalFile("./static", true)))

	store := cookie.NewStore([]byte(utils.AppConfig.Server.SecretKey))
	r.Use(sessions.Sessions("Status", store))

	r.LoadHTMLGlob("static/*/*.tmpl")

	html := r.Group("/")
	html.Use(middleware.Login())
	html.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.tmpl", nil) })

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/test", handler.GenAuthTokenHandler)

	r.GET("/login", handler.LoginHandler)
	r.GET("/login/:app", handler.LoginHandler)
	r.GET("/logout", handler.LogoutHandler)
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
		login.GET("/decode", handler.DecodeHandler)
	}

	r.GET("/auth", handler.AuthHandler)
	r.GET("/auth/:app", handler.AuthAPPHandler)

	weixin := r.Group("/weixin")
	weixin.Use(middleware.Auth())
	{
		weixin.GET("/department", handler.GetDepartmentListHandler)
		weixin.GET("/department/:departmentID", handler.GetDepartmentUsersHandler)
		weixin.PATCH("/user/:userID", handler.UpdateUserInfoHandler)
	}

	Sms := r.Group("/sms")
	Sms.Use(middleware.Auth())
	{
		Sms.POST("/send_single", handler.SendSingleSMSHandler)
		Sms.POST("/send_group", handler.SendGroupSMSHandle)
		Sms.GET("/templates", handler.GetTemplatesHandler)
		Sms.GET("/reply_callback", handler.ReplyCallbackHandler)
	}

	open := r.Group("/open")
	open.Use(middleware.Auth())
	{
		open.GET("/permission", handler.GetPermissionHandler)
	}

	mail := r.Group("/mail")
	mail.Use(middleware.Auth())
	{
		mail.POST("/send_mail", handler.SendMailHandler)
	}

	r.GET("/v/:Shorturl", handler.MapShortUrlHandler)
	r.GET("/genUrl", middleware.Auth(), handler.CreateShortUrlHandler)

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
