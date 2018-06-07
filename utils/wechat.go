package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"open-platform/wechat"
	"strconv"
	"time"
)

// Contact status work weixin api
var Contact wechat.WorkWeixin
var Login wechat.WorkWeixin

type codeResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	UserID  string `json:"UserId"`
}

// UserInfo handle user info
type UserInfo struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	Name        string `json:"name"`
	UserID      string `json:"userid"`
	Departments []int  `json:"department"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	IsLeader    int    `json:"isleader"`
}

func init() {
	Contact.Init(AppConfig.WeWork.CropID, AppConfig.WeWork.ContactSecret, 0)
	Login.Init(AppConfig.WeWork.CropID, AppConfig.WeWork.AgentSecret, AppConfig.WeWork.AgentID)
}

// MakePCRedirctString to make up redirect auth string
func MakePCRedirctString(redirectURL, state string) string {
	u := url.Values{}
	u.Set("appid", AppConfig.WeWork.CropID)
	u.Set("agentid", strconv.Itoa(AppConfig.WeWork.AgentID))
	u.Set("redirect_uri", redirectURL)
	u.Set("state", state)

	reuqestURL := "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?" + u.Encode()
	return reuqestURL
}

// MakeMobileRedirctString to make up redirect auth string
func MakeMobileRedirctString(redirectURL, state, scope string) string {
	u := url.Values{}
	u.Set("response_type", "code")
	u.Set("scope", scope)
	u.Set("appid", AppConfig.WeWork.CropID)
	// scope
	// snsapi_base：静默授权，可获取成员的基本信息；
	// snsapi_UserInfo：静默授权，可获取成员的敏感信息，但不包含手机、邮箱；
	// snsapi_privateinfo：手动授权，可获取成员的敏感信息，包含手机、邮箱
	u.Set("agentid", strconv.Itoa(AppConfig.WeWork.AgentID))
	u.Set("redirect_uri", redirectURL)
	u.Set("state", state)

	reuqestURL := "https://open.weixin.qq.com/connect/oauth2/authorize?" + u.Encode() + "#wechat_redirect"

	return reuqestURL
}

// VerifyCode is a func to verify work wexin code
func VerifyCode(code string) (UserID string, Error error) {

	u := url.Values{}
	u.Set("access_token", Login.GetAccessToken())
	u.Set("code", code)
	reuqestURL := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?" + u.Encode()

	fmt.Println(reuqestURL)

	data := new(codeResp)
	getJSON(reuqestURL, "", data, nil)
	fmt.Println("data", data)

	if data.ErrCode != 0 {
		return data.ErrMsg, errors.New(data.ErrMsg)
	}
	return data.UserID, nil
}

// GetUserInfo is a func to get user info
func GetUserInfo(userid string) (data *UserInfo, Error error) {
	u := url.Values{}
	u.Set("access_token", Contact.GetAccessToken())
	u.Set("userid", userid)
	reuqestURL := "https://qyapi.weixin.qq.com/cgi-bin/user/get?" + u.Encode()

	info := new(UserInfo)
	getJSON(reuqestURL, "", info, nil)

	if info.ErrCode == 0 {
		return info, nil
	}
	return info, errors.New(info.ErrMsg)
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, header string, target, data interface{}) error {
	// myClient.Header.Set()
	// r, err := myClient.Get(url)

	req, err := http.NewRequest("GET", url, nil)
	// ...
	if header != "" {
		req.Header.Add("Authorization", header)
	}
	r, err := myClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
