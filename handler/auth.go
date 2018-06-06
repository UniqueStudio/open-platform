package handler

import (
	"net/http"
	"net/url"
	"open-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthHandler is a func to resolve auth third party requests
func AuthHandler(c *gin.Context) {
	state := utils.B64Decode(c.Query("state"))
	code := c.Query("code")
	userID, _ := utils.VerifyCode(code)
	userInfo, _ := utils.GetUserInfo(userID)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "info": userInfo, "state": state})
}

// TestHandler is a func to resolve auth third party requests
func TestHandler(c *gin.Context) {
	u := url.Values{}
	u.Set("appid", utils.AppConfig.WeWork.CropID)
	u.Set("agentid", strconv.Itoa(utils.AppConfig.WeWork.AgentID))
	u.Set("redirect_uri", "https://test.fredliang.cn/auth")
	u.Set("state", "state")

	reuqestURL := "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?" + u.Encode()
	c.Redirect(http.StatusFound, reuqestURL)
}
