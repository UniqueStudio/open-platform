package handler

import (
	"net/http"
	"open-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

// LoginHandler is a func to handle user login request
func LoginHandler(c *gin.Context) {
	UserAgent := c.GetHeader("User-Agent")
	state := c.Query("state")
	ua := user_agent.New(UserAgent)

	if ua.Mobile() {
		app := c.Param("app")
		c.Redirect(http.StatusFound, utils.MakeMobileRedirctString("https://"+utils.AppConfig.Server.Hostname+"auth/"+app, state, "snsapi_userinfo"))
	} else {
		c.HTML(http.StatusOK, "html/index_pc.tmpl", gin.H{
			"redirURL": utils.MakePCRedirctString("https://"+utils.AppConfig.Server.Hostname+"/check", state),
		})
	}
}
