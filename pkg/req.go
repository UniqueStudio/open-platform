package pkg

type SingleSMSReq struct {
	PhoneNumber      string   `json:"phone_number"`
	TemplateParamSet []string `json:"template_param_set"`
	TemplateID       uint     `json:"template_id"`
	SignID           *uint    `json:"sign_id,omitempty"`
}

type MultipleSMSReq struct {
	PhoneNumber      []string `json:"phone_number"`
	TemplateParamSet []string `json:"template_param_set"`
	TemplateID       uint     `json:"template_id"`
	SignID           *uint    `json:"sign_id,omitempty"`
}

type AddSignReq struct {
	SignID string `json:"sign_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type AddTemplateReq struct {
	TemplateID  string `json:"template_id" binding:"required"`
	Content     string `json:"content" binding:"required"`
	ParamNumber int32  `json:"param_number" binding:"required"`
}

type AddEmailTemplateReq struct {
	TemplateName string `json:"template_name" binding:"required"`
	Content      string `json:"content" binding:"required"`
	ParamNumber  int32  `json:"param_number" binding:"required"`
}

type SingleEmailReq struct {
	EmailTo          string   `json:"email" binding:"required"`
	TemplateParamSet []string `json:"template_param_set" binding:"required"`
	TemplateID       uint     `json:"template_id" binding:"required"`
	Subject          string   `json:"subject" binding:"required"`
}

type MultipleEmailReq struct {
	EmailTo          []string   `json:"email" binding:"required"`
	TemplateParamSet [][]string `json:"template_param_set" binding:"required"`
	TemplateID       uint       `json:"template_id" binding:"required"`
	Subject          string     `json:"subject" binding:"required"`
}

type EmailTemplatesReq struct {
	Content     string `json:"content" binding:"required"`
	TemplateID  uint32 `json:"template_id" binding:"required"`
	ParamNumber uint32 `json:"param_number"`
}
