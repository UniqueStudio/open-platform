package db

type Reply struct {
	Id         int64  `xorm:"pk autoincr"`
	Extend     string `json:"extend"`
	Mobile     string `json:"mobile"`
	NationCode string `json:"nationcode"`
	Sign       string `json:"sign"`
	Text       string `json:"text"`
	Time       int    `json:"time"`
}

type Status struct {
	Id              int64  `xorm:"pk autoincr"`
	UserReceiveTime string `json:"user_receive_time"`
	NationCode      string `json:"nationcode"`
	Mobile          string `json:"mobile"`
	ReportStatus    string `json:"report_status"`
	Errmsg          string `json:"errmsg"`
	Description     string `json:"description"`
	Sid             string `json:"sid"`
}
