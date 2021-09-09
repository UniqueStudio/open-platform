package pkg

type SingleSMSReq struct {
	PhoneNumber      string   `json:"phone_number"`
	TemplateParamSet []string `json:"template_param_set"`
	TemplateID       string   `json:"template_id"`
}

type MultipleSMSReq struct {
	MultipaleSMS []SingleSMSReq `json:"multipale_sms"`
}
