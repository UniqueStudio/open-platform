package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jochasinga/requests"
	"net/http"
	"net/url"
	"open-platform/wechat"
	"strconv"
	"time"
)

// Contact status work weixin api
var Contact wechat.WorkWeixin

// Login status work weixin api
var Login wechat.WorkWeixin

type codeResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	UserID  string `json:"UserId"`
}

// UserInfo handle user info
type UserInfo struct {
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
	UserID     string `json:"userid"`
	Name       string `json:"name,omitempty"`
	Department []int  `json:"department,omitempty"`
	Order      []int  `json:"order,omitempty"`
	Position   string `json:"position,omitempty"`
	Mobile     string `json:"mobile,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Email      string `json:"email,omitempty"`
	IsLeader   int    `json:"isleader,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
	Telephone  string `json:"telephone,omitempty"`
	Enable     int    `json:"enable,omitempty"`
	Alias      string `json:"alias,omitempty"`
	Extattr    struct {
		Attrs []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"attrs"`
	} `json:"extattr,omitempty"`
	Status           int    `json:"status,omitempty"`
	QrCode           string `json:"qr_code,omitempty"`
	ExternalPosition string `json:"external_position,omitempty"`
	ExternalProfile  struct {
		ExternalAttr []struct {
			Type int    `json:"type"`
			Name string `json:"name"`
			Text struct {
				Value string `json:"value"`
			} `json:"text,omitempty"`
			Web struct {
				URL   string `json:"url"`
				Title string `json:"title"`
			} `json:"web,omitempty"`
			Miniprogram struct {
				Appid    string `json:"appid"`
				Pagepath string `json:"pagepath"`
				Title    string `json:"title"`
			} `json:"miniprogram,omitempty"`
		} `json:"external_attr,omitempty"`
	} `json:"external_profile,omitempty"`
}

// UpdateInfoStruct is a type for udate service 
type UpdateInfoStruct struct {
	Userid        string `json:"userid"`
	Name          string `json:"name,omitempty"`
	Department    []int  `json:"department,omitempty"`
	Order         []int  `json:"order,omitempty"`
	Position      string `json:"position,omitempty"`
	Mobile        string `json:"mobile,omitempty"`
	Gender        string `json:"gender,omitempty"`
	Email         string `json:"email,omitempty"`
	Isleader      int    `json:"isleader,omitempty"`
	Enable        int    `json:"enable,omitempty"`
	AvatarMediaid string `json:"avatar_mediaid,omitempty"`
	Telephone     string `json:"telephone,omitempty"`
	Alias         string `json:"alias,omitempty"`
	Extattr       struct {
		Attrs []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"attrs,omitempty"`
	} `json:"extattr,omitempty"`
	ExternalPosition string `json:"external_position,omitempty"`
	ExternalProfile  struct {
		ExternalAttr []struct {
			Type int    `json:"type"`
			Name string `json:"name"`
			Text struct {
				Value string `json:"value"`
			} `json:"text,omitempty"`
			Web struct {
				URL   string `json:"url"`
				Title string `json:"title"`
			} `json:"web,omitempty"`
			Miniprogram struct {
				Appid    string `json:"appid"`
				Pagepath string `json:"pagepath"`
				Title    string `json:"title"`
			} `json:"miniprogram,omitempty"`
		} `json:"external_attr,omitempty"`
	} `json:"external_profile,omitempty"`
}

func init() {
	Contact.Init(AppConfig.WeWork.CropID, AppConfig.WeWork.ContactSecret, 0)
	Login.Init(AppConfig.WeWork.CropID, AppConfig.WeWork.AgentSecret, AppConfig.WeWork.AgentID)
}

// GetDepartmentUsers is a func to get group info
func GetDepartmentUsers(groupsID []int) (data []wechat.User, Error error) {

	fmt.Println("groupsID:", groupsID)
	var userList []wechat.User

	for _, departmentID := range groupsID {
		fmt.Println("departmentID: ", departmentID)

		tempList := Contact.GetDepartmentUsers(departmentID, 1)
		fmt.Println("tempList: ", tempList)

		userList = append(userList, tempList...)

	}
	fmt.Println("return userList")
	return userList, nil
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

// UpdateUserInfo is a func to update user info
func UpdateUserInfo(data UpdateInfoStruct) (Error error) {
	u := url.Values{}
	u.Set("access_token", Contact.GetAccessToken())
	reuqestURL := "https://qyapi.weixin.qq.com/cgi-bin/user/update?" + u.Encode()

	res, err := requests.PostJSON(reuqestURL, data)
	fmt.Println(reuqestURL, res.String())

	if err != nil {
		return err
	}

	var resp codeResp
	json.Unmarshal(res.JSON(), &resp)

	if resp.ErrCode != 0 {
		return errors.New(string(res.JSON()))
	}
	fmt.Println(resp.ErrCode)

	return nil
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
