package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"open-platform/utils"
	"strconv"
	"time"
)

// AuthHandler is a func auth user request
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

// AuthAPPHandler is a func auth user request
func AuthAPPHandler(c *gin.Context) {
	session := sessions.Default(c)

	code := c.Query("code")
	state := c.Query("state")
	appID := c.Query("appid")
	app := c.Param("app")

	switch state {
	case "":
		fmt.Println(code)
		UserID, err := utils.VerifyCode(code)

		info, _ := utils.GetUserInfo(UserID)
		fmt.Println("info", info)

		session.Set("UserID", info.UserID)
		session.Set("IsLeader", info.IsLeader == 1)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"UserId":   UserID,
			"phone":    info.Mobile,
			"username": info.Name,
			"state":    state,
			"appid":    appID,
			"err":      err})

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
		c.Redirect(http.StatusFound, "app/"+app+"?code="+code+"&state="+state)
	}
}

// DecodeHandler is a func to resolve auth third party requests
func DecodeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"state": utils.B64Decode(c.Query("state")), "code": http.StatusOK})
}
