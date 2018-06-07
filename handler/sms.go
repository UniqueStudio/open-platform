package handler

import (
	"net/http"
	"open-platform/utils"

	"github.com/gin-gonic/gin"
)

// GetSMSTemplate is a func to get sms template
func GetSMSTemplate(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{"code":http.StatusOK,"data":utils.GetQCSMSTemplate().Data})
}