package handler

import (
	"fmt"
	"net/http"
	"open-platform/utils"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetDepartmentUsersHandler is a func to get sms template
func GetDepartmentUsersHandler(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("departmentID"))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": err, "code": http.StatusConflict})
	}

	userInfo, err := utils.GetDepartmentUsers([]int{departmentID})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": err, "code": http.StatusConflict})
	}

	c.JSON(http.StatusOK, gin.H{"message": err, "code": http.StatusOK, "data": userInfo})
}

// GetDepartmentListHandler is a func to get wexin work department list
func GetDepartmentListHandler(c *gin.Context) {
	data, err := utils.Contact.GetDepartmentParentList()

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": err, "code": http.StatusConflict})
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": data, "code": http.StatusOK})
}

// CheckAuthorityHandler is a func to check user Authority
func CheckAuthorityHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	appid := c.Query("appid")
	session := sessions.Default(c)

	switch state {
	case "test":
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
				"message": err,
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
