package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Token func is a middleware to check if user is logged in
func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		UserID := session.Get("UserID")
		IsLeader := session.Get("IsLeader")

		fmt.Println("UserID, IsLeader", UserID, IsLeader)
		if UserID == nil {
			state := string([]byte(c.Request.URL.Path)[1:])
			c.Redirect(http.StatusFound, "/login?state="+state)
			c.Abort()
		} else {
			c.Set("UserID", UserID)
			c.Set("IsLeader", IsLeader)
			c.Next()
		}
	}
}
