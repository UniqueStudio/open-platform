package middleware

import (
	"fmt"
	"net/http"
	"open-platform/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Auth is a middleware to verify access
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		Token := c.GetHeader("Token")
		session := sessions.Default(c)
		LoginUserID := session.Get("UserID")
		IsLeader := session.Get("IsLeader")

		fmt.Println("Token: ", Token)
		switch Token {
		case "":
			if LoginUserID != nil {
				c.Set("UserID", LoginUserID)
				c.Set("IsLeader", IsLeader)
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Empty Token Please authorize before requesting"})
				c.Abort()
			}

		default:
			UserID, IsLeader, err := utils.LoadToken(Token)

			if LoginUserID != nil {
				c.Set("UserID", UserID)
				c.Set("IsLeader", IsLeader)
				c.Next()
			} else {
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"message": "Please authorize before requesting"})
					c.Abort()
				}

				c.Set("UserID", UserID)
				c.Set("IsLeader", IsLeader)

				c.Next()
			}
		}
	}
}

// Admin is a middleware to check admin
// You Have to append Header:
// `Token: oausudgaosugdoa``
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		Token := c.GetHeader("Token")

		switch Token {
		case "":
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Please authorize before requesting"})
			c.Abort()

		default:
			UserID, IsLeader, err := utils.LoadToken(Token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Please authorize before requesting"})
				c.Abort()
			}
			if IsLeader == false {
				c.JSON(http.StatusForbidden, gin.H{"message": "Please use an admin token to request"})
				c.Abort()
			}

			c.Set("UserID", UserID)
			c.Set("IsLeader", IsLeader)

			c.Next()
		}
	}
}

//Login func is used to check if user is logged in
func Login() gin.HandlerFunc {
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
