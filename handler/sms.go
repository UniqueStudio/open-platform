package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"open-platform/db"
	"open-platform/utils"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// 请求单次发送的json格式
type smsSendInfo struct {
	PhoneNumber      string   `json:"phone_number"`
	TemplateParamSet []string `json:"template_param_set"`
	TemplateID       string   `json:"template_id"`
}

func SendSingleSMSHandler(c *gin.Context) {
	var data smsSendInfo
	c.BindJSON(&data)
	log.Println(data)
	err := CheckSMS(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "data": err.Error()})
		return
	}
	data.PhoneNumber = "+86" + data.PhoneNumber
	smsResponse, err := utils.SendSingleSms(data.PhoneNumber, data.TemplateParamSet, data.TemplateID)
	if err != nil {
		// 应该是服务端的问题
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "data": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": smsResponse})
}

func SendGroupSMSHandle(c *gin.Context) {
	var wg sync.WaitGroup
	var data []smsSendInfo
	c.BindJSON(&data)
	log.Println(data)
	var errMsgs []string
	var responses []*sms.SendSmsResponse
	for _, perSms := range data {
		err := CheckSMS(&perSms)
		if err != nil {
			errMsgs = append(errMsgs, err.Error())
			continue
		}
		perSms.PhoneNumber = "+86" + perSms.PhoneNumber
		wg.Add(1)
		go func(perSms smsSendInfo) {
			defer wg.Done()
			smsResponse, err := utils.SendSingleSms(perSms.PhoneNumber, perSms.TemplateParamSet, perSms.TemplateID)
			if err != nil {
				errMsgs = append(errMsgs, err.Error())
			}
			responses = append(responses, smsResponse)
		}(perSms)
	}
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": responses, "error": errMsgs})
}

func GetTemplatesHandler(c *gin.Context) {
	allTemplates := utils.GetTemplates()
	if len(allTemplates) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": http.NotFound, "data": "No templates found."})
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": allTemplates})
}

func CheckSMS(data *smsSendInfo) error {
	// 模板存在
	allTemplates := utils.GetTemplates()
	var template *utils.SMSTemplate = nil
	for _, temp := range allTemplates {
		if temp.ID == data.TemplateID {
			template = temp
			break
		}
	}
	paramSetString := strings.Join(data.TemplateParamSet, ",")
	if template == nil {
		return errors.New(fmt.Sprintf("[Paraments:%s] [TemplateID:%s]:Template not exists", paramSetString, data.TemplateID))
	}
	// 参数匹配
	if template.ParamNum != len(data.TemplateParamSet) {
		return errors.New(fmt.Sprintf("[Paraments:%s] ParamNumber error. Template Content:%s",
			paramSetString, template.Content))
	}
	// 电话号码，11位且数字（不加+86）
	match, _ := regexp.Match(`^[0-9]{11}$`, []byte(data.PhoneNumber))
	if !match {
		return errors.New(fmt.Sprintf("[Paraments:%s] [PhoneNumber:%s] Format error", paramSetString, data.PhoneNumber))
	}
	return nil
}

func ReplyCallbackHandler(c *gin.Context) {
	fmt.Println(c.Params)
	var statusData []db.Status
	c.BindJSON(&statusData)

	_, err := db.ORM.Insert(&statusData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK"})
}
