package utils

// StdResp is an example respond struct
type StdResp struct {
	Message string `json:"message" example:"OK"`
	Code    int    `json:"code" example:"200"`
}

type SMSTemplateResp struct {
	Message string `json:"message" example:"OK"`
	Code    int    `json:"code" example:"200"`
	Data    struct {
		ID            uint   `json:"id" example:"96385"`
		Text          string `json:"text" example:"您注册{1}的验证码是{2}"`
		Status        uint   `json:"status" example:"0"`
		Reply         string `json:"reply" example:""`
		Type          uint   `json:"type" example:"0"`
		International uint   `json:"international" example:"0"`
		ApplyTime     string `json:"apply_time" example:"2018-03-18 15:41:18"`
	} `json:"data"`
}

type DepartmentUsersResp struct {
	Message string `json:"message" example:"OK"`
	Code    int    `json:"code" example:"200"`
	Data    struct {
		Userid       string `json:"userid" example:"userid"`
		Name         string `json:"name" example:"name"`
		Department   []int  `json:"department" example:"2,3"`
		Mobile       string `json:"mobile" example:"17371266666"`
		Email        string `json:"email" example:"ahsudhoa@jzdbcadg"`
		Status       int    `json:"status" example:"2"`
		Avatar       string `json:"avatar" example:"http://p.qlogo.cn/bizmail/LIdibicNn9rcMNTXq4HzI8vkYib9XvU4H1mTgIonBt5gy4ibLNtuu"`
		Telephone    string `json:"telephone" example:""`
		English_name string `json:"english_name" example:"fred"`
	} `json:"data"`
}




type DepartmentListResp struct {
	Message string `json:"message" example:"OK"`
	Code    int    `json:"code" example:"200"`
	Data    struct {
		Id       int    `json:"id" example:"2"`
		Name     string `json:"name" example:"Design"`
		Parentid int32  `json:"parentid" example:"3"`
		Order    int32  `json:"order" example:"99999000"`
	} `json:"data"`
}

// ErrResp is an example respond struct
type ErrResp struct {
	Message string `json:"message" example:"Error Message"`
	Code    int    `json:"code" example:"400"`
}
