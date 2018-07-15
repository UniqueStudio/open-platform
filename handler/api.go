package handler

import (
	"net/http"
	"open-platform/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GenAccessKeyHandler is a func generate access key
func GenAccessKeyHandler(c *gin.Context) {
	session := sessions.Default(c)
	state := c.Query("state")

	UserID := session.Get("UserID")
	if UserID == nil {
		c.Redirect(http.StatusFound, "/login?state=api")
		return
	}
	//IsLeader :=session.Get("IsLeader")

	switch state {
	case "":
		c.Redirect(http.StatusFound, "/login?state=api")

	case "api":

		if UserID == "" {
			c.Redirect(http.StatusFound, "/login?state=api")
		} else {
			info, _ := utils.GetUserInfo(UserID.(string))
			c.JSON(http.StatusOK, gin.H{"accessKey": utils.GenAccessKey(info.UserID, info.IsLeader == 1), "status": "OK"})
		}
	}
}
