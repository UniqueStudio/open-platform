package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderAppStaticFilesHandler is a func to render static files
func RenderAppStaticFilesHandler(c *gin.Context) {
	app := c.Param("app")
	isLeader, exist := c.Get("isAdmin")
	if !isLeader.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not admin", "code": http.StatusForbidden})
	}
	if !exist {
		c.Redirect(http.StatusFound, "./login/"+app)
	}
	c.HTML(http.StatusOK, app+"/dist/index.html", nil)
}
