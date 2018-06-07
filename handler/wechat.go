package handler

import (
	"fmt"
	"net/http"
	"open-platform/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CheckAuthorityHandler is a func to check Authority
func CheckAuthorityHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	appid := c.Query("appid")
	session := sessions.Default(c)

	switch state {
	case "":
		fmt.Println(code)
		UserID, err := utils.VerifyCode(code)

		info, _ := utils.GetUserInfo(UserID)
		fmt.Println("info", info)

		session.Set("UserID", info.UserID)
		session.Set("IsLeader", info.IsLeader == 1)
		session.Save()

		c.JSON(http.StatusOK, gin.H{"UserId": UserID,
			"phone": info.Mobile, "username": info.Name,
			"state": state, "appid": appid, "err": err})

	default:
		UserID, err := utils.VerifyCode(code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err,
			})
		}
		info, _ := utils.GetUserInfo(UserID)
		fmt.Println("info", info)

		fmt.Println("case test:")
		session.Set("UserID", info.UserID)
		session.Set("IsLeader", info.IsLeader == 1)
		session.Save()

		fmt.Println("state:", state)
		c.Redirect(http.StatusFound, "/"+state+"?code="+code+"&state="+state)

	}
}
