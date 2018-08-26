package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"open-platform/utils"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userPermisson struct {
	Avatar      string   `json:"avatar"`
	UserID      string   `json:"userID"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// GetPermissionHandler is a func to get user permission
func GetPermissionHandler(c *gin.Context) {
	session := sessions.Default(c)

	tmp := session.Get("UserID")

	if tmp == nil {
		c.JSON(http.StatusNonAuthoritativeInfo, gin.H{"message": "Empty UserID", "code": http.StatusNonAuthoritativeInfo})
		return
	}

	userID := tmp.(string)

	permisson := userPermisson{UserID: userID}
	userInfo, err := utils.GetUserInfo(userID)
	permisson.Avatar = userInfo.Avatar
	if userInfo.IsLeader == 1 {
		permisson.Permissions = []string{"admin", "write", "read"}
	} else {
		permisson.Permissions = []string{"read"}
	}

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": userInfo.ErrMsg, "code": http.StatusConflict})
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "code": http.StatusOK, "data": permisson})
}

// AuthHandler is a func auth user request
func AuthHandler(c *gin.Context) {
	var state utils.State
	decoded, err := utils.B64Decode(c.Query("state"))
	err = json.Unmarshal([]byte(decoded), &state)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"code": http.StatusConflict, "message": err})
		return
	}

	AccessKey := state.AccessKey
	if AccessKey == "" {
		AccessKey = state.Token
	}

	APPUserID, _, err := utils.LoadAccessKey(AccessKey)

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

	case "api":
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
		c.Redirect(http.StatusFound, "/api?code="+code+"&state="+state)

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
		c.Redirect(http.StatusFound, "/app/"+app+"?code="+code+"&state="+state)
	}
}

// DecodeHandler is a func to resolve auth third party requests
func DecodeHandler(c *gin.Context) {
	decoded, err := utils.B64Decode(c.Query("state"))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"code": http.StatusConflict, "message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"state": decoded, "code": http.StatusOK})
}

// GenAuthTokenHandler is a func to gen Auth Token
func GenAuthTokenHandler(c *gin.Context) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"expire": func() int64 {
			now := time.Now()
			duration, _ := time.ParseDuration("14d")
			m1 := now.Add(duration)
			return m1.Unix()
		}(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(utils.AppConfig.Server.SecretKey))

	fmt.Println(tokenString, err)
	c.String(http.StatusOK, tokenString)
}
