package pkg

import (
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type UniformResource struct {
	Code    int         `json:"code"`
	State   bool        `json:"state,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ErrorResponse(err error) *UniformResource {
	return &UniformResource{
		State:   false,
		Message: err.Error(),
		Data:    err,
	}
}

func SuccessResponse(data interface{}) *UniformResource {
	return &UniformResource{
		Code:    200,
		State:   true,
		Message: "ok",
		Data:    data,
	}
}

type SMSSendStatus struct {
	SerialNo    string `json:"serial_no,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Fee         uint64 `json:"fee,omitempty"`
	ErrCode     string `json:"err_code,omitempty"`
	Message     string `json:"message,omitempty"`
}

type EmailSendStatus struct {
	Err     string `json:"err_code,omitempty"`
	To      string `json:"to,omitempty"`
	Content string `json:"content,comitempty"`
}

func TencentSMSToSMSResp(resp *sms.SendSmsResponse) *[]SMSSendStatus {
	if resp.Response == nil || resp.Response.SendStatusSet == nil {
		return nil
	}
	statuses := make([]SMSSendStatus, len(resp.Response.SendStatusSet))
	for i := range resp.Response.SendStatusSet {
		status := SMSSendStatus{}
		if resp.Response.SendStatusSet[i].SerialNo != nil {
			status.SerialNo = *resp.Response.SendStatusSet[i].SerialNo
		}
		if resp.Response.SendStatusSet[i].PhoneNumber != nil {
			status.PhoneNumber = *resp.Response.SendStatusSet[i].PhoneNumber
		}
		if resp.Response.SendStatusSet[i].Fee != nil {
			status.Fee = *resp.Response.SendStatusSet[i].Fee
		}
		if resp.Response.SendStatusSet[i].Code != nil {
			status.ErrCode = *resp.Response.SendStatusSet[i].Code
		}
		if resp.Response.SendStatusSet[i].Message != nil {
			status.Message = *resp.Response.SendStatusSet[i].Message
		}
		statuses[i] = status
	}
	return &statuses
}
