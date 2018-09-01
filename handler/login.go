package handler

import (
	"net/http"
	"open-platform/utils"
	"strconv"

	"github.com/gin-contrib/sessions"
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
		c.Redirect(http.StatusTemporaryRedirect, utils.MakeMobileRedirctString("https://"+utils.AppConfig.Server.Hostname+"/auth/"+app, state, "snsapi_userinfo"))
	} else {
		c.HTML(http.StatusOK, "index_pc.tmpl", gin.H{
			"redirURL":     utils.MakePCRedirctString("https://"+utils.AppConfig.Server.Hostname+"/check", state),
			"appid":        utils.AppConfig.WeWork.CropID,
			"agentid":      strconv.Itoa(utils.AppConfig.WeWork.AgentID),
			"redirect_uri": "https%3a%2f%2fopen.hustunique.com%2fcheck",
			"state":        state,
		})
	}
}

// LogoutHandler is a func to handle user login request
func LogoutHandler(c *gin.Context) {
	state := c.Query("state")

	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/"+state)

}
