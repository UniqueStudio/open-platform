package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"open-platform/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandler is a func to resolve auth third party requests
func AuthHandler(c *gin.Context) {
	var state utils.State
	err := json.Unmarshal([]byte(utils.B64Decode(c.Query("state"))), &state)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"code": http.StatusConflict, "message": err})
		return
	}
	APPUserID, _, err := utils.LoadToken(state.Token)
	if err != nil || APPUserID == "" {
		c.JSON(http.StatusNonAuthoritativeInfo, gin.H{"code": http.StatusNonAuthoritativeInfo, "message": err})
		return
	}

	fmt.Println("redirURL", state.URL)

	code := c.Query("code")
	userID, _ := utils.VerifyCode(code)
	userInfo, _ := utils.GetUserInfo(userID)

	u := url.Values{}
	data, _ := json.Marshal(userInfo)
	u.Set("state", utils.B64Encode(string(data)))
	u.Set("timestamp", fmt.Sprintln(time.Now().Unix()))
	c.Redirect(http.StatusFound, state.URL+"?"+u.Encode())
}

// TestHandler is a func to resolve auth third party requests
func TestHandler(c *gin.Context) {
	u := url.Values{}
	u.Set("appid", utils.AppConfig.WeWork.CropID)
	u.Set("agentid", strconv.Itoa(utils.AppConfig.WeWork.AgentID))
	u.Set("redirect_uri", "https://test.fredliang.cn/auth")
	state := utils.B64Encode(`{"url":"https://test.fredliang.cn/decode"}`)

	fmt.Println("state", state)
	u.Set("state", state)

	reuqestURL := "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?" + u.Encode()
	c.Redirect(http.StatusFound, reuqestURL)
}

// DecodeHandler is a func to resolve auth third party requests
func DecodeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"state": utils.B64Decode(c.Query("state")), "code": http.StatusOK})
}
