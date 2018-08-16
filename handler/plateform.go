package handler

import (
	"net/http"
	"open-platform/utils"

	"github.com/gin-gonic/gin"
)

type platformSMSRequest struct {
	Template     int      `json:"template"`
	ParamList    []string `json:"param_list"`
	GroupChosen  []int    `json:"group_chosen"`
	MemberChosen []string `json:"member_chosen"`
}

// PlatformSendSMS is a func to send sms
func PlatformSendSMS(c *gin.Context) {
	var data platformSMSRequest
	c.BindJSON(&data)
	if data.Template == 0 || (len(data.GroupChosen)+len(data.MemberChosen) == 0) {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Missing parameter", "code": http.StatusNotAcceptable})
	}

	phoneList := []string{}

	// get phone numbers from member list
	for _, userID := range data.MemberChosen {
		userInfo, err := utils.GetUserInfo(userID)
		phoneList = append(phoneList, userInfo.Mobile)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error(), "code": http.StatusConflict})
			return
		}
	}

	// get phone numbers from group list
	userList, err := utils.GetDepartmentUsers(data.GroupChosen)
	for _, user := range userList {
		phoneList = append(phoneList, user.Mobile)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error(), "code": http.StatusConflict})
			return
		}
	}

	phoneList = utils.RemoveDuplicate(phoneList)

	// send sms
	isOK, message, errID := utils.SendQCSMSMulti(phoneList, data.Template, data.ParamList)

	if isOK {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": message, "error_id": errID})

}
