package database

import "gorm.io/gorm"

type SMSSignature struct {
	gorm.Model `json:"-"`
	VirtualID  string `gorm:"virtual_id"`
	ID         string `gorm:"id"`
	Name       string `gorm:"name"`
}

type SMSTemplate struct {
	gorm.Model  `json:"-"`
	VirtualID   string `gorm:"virtual_id"`
	ID          string `gorm:"id"`
	Content     string `gorm:"content"`
	ParamNumber int    `gorm:"param_number"`
}

func (sign *SMSSignature) TableName() string {
	return "sms_signature"
}

func (template *SMSTemplate) TableName() string {
	return "sms_template"
}
