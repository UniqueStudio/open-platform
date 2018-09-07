package handler

import (
	"net/http"
	"open-platform/db"
	"open-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type sms struct {
	Phone     string   `json:"phone"`
	Template  int      `json:"template"`
	ParamList []string `json:"param_list"`
}

// GetSMSTemplateHandler is a func to get sms template
func GetSMSTemplateHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": utils.GetQCSMSTemplate().Data})
}

// GetSMSTemplateStatusHandler is a func to get sms template
func GetSMSTemplateStatusHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	data, err := utils.GetQCSMSTemplateStatus([]uint{uint(id)})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})

	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": data})
}

// AddSMSTemplateHandler is a func to get sms template
func AddSMSTemplateHandler(c *gin.Context) {
	var info templateInfo
	c.BindJSON(&info)
	data, err := utils.AddQCSMSTemplate(info.Title, info.Text, info.Remark)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": data})
}

// SendSMSHandler is a func to send sms via sms template
func SendSMSHandler(c *gin.Context) {
	var data sms
	c.BindJSON(&data)
	isOK, message, errID := utils.SendQCSMS(data.Phone, data.Template, data.ParamList)

	if isOK {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": message, "error_id": errID})

}

// ReplyCallbackHandler is a handler to receive sms reply callback request
func ReplyCallbackHandler(c *gin.Context) {
	var replyData db.Reply
	c.BindJSON(&replyData)
	_, err := db.ORM.Insert(&replyData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK"})

}

// StatusCallbackHandler is a handler to receive sms Status callback request
func StatusCallbackHandler(c *gin.Context) {
	var statusData []db.Status
	c.BindJSON(&statusData)

	_, err := db.ORM.Insert(&statusData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK"})
}

// GetReplyHandler is a handler to receive sms Status callback request
func GetReplyHandler(c *gin.Context) {
	replyList := make([]db.Reply, 0)

	err := db.ORM.Find(&replyList)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": replyList})
}

// GetStatusHandler is a handler to receive sms Status callback request
func GetStatusHandler(c *gin.Context) {
	statusList := make([]db.Status, 0)

	err := db.ORM.Find(&statusList)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": statusList})
}

type templateInfo struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Remark string `json:"remark"`
}
