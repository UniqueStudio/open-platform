package db

type Short_Url struct {
	Id       int64  `json:"id" xorm:"pk autoincr"`
	Url      string `json:"url"`
	Shorturl string `json:"shorturl"`
	Hashcode string `json:"hashcode"`
}
