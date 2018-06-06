package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"open-platform/wechat"
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
