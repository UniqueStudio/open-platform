package handler

import (
	"net/http"
	"open-platform/utils"

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
