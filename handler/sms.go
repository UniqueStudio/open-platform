package handler

import (
	"net/http"
	"open-platform/utils"

	"github.com/gin-gonic/gin"
)

// GetSMSTemplate is a func to get sms template
// @Summary Get SMS Template
// @Description Get SMS Template
// @Tags sms
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} utils.SMSTemplateResp
// @Router /weixin/template [get]
func GetSMSTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": utils.GetQCSMSTemplate().Data})
}
