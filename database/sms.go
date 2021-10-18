package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SMSSignature struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ID     uint   `gorm:"column:virtual_id;primarykey" json:"template_id"`
	SignID string `gorm:"column:sign_id" json:"-"`
	Name   string `gorm:"column:name" json:"name"`
}

type SMSTemplate struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ID          uint   `gorm:"column:virtual_id;primarykey" json:"template_id"`
	TemplateID  string `gorm:"column:template_id" json:"-"`
	Content     string `gorm:"column:content" json:"content"`
	ParamNumber int32  `gorm:"column:param_number" json:"param_number"`
}

func (sign *SMSSignature) TableName() string {
	return "sms_signature"
}

func (template *SMSTemplate) TableName() string {
	return "sms_template"
}

func InsertSMSSignatures(signs *[]SMSSignature) (*[]SMSSignature, error) {
	result := OpenDB.Create(signs)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, errors.New("no data inserted")
	}
	return signs, nil
}

func GetSMSSignatureByVirtualID(ctx context.Context, virtualId uint) (*SMSSignature, error) {
	sign := new(SMSSignature)
	result := OpenDB.WithContext(ctx).Table(sign.TableName()).Where("virtual_id = ?", virtualId).Scan(sign)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, gorm.ErrRecordNotFound
	}
	return sign, nil
}

func InsertSMSTemplates(templates *[]SMSTemplate) (*[]SMSTemplate, error) {
	result := OpenDB.Create(templates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, errors.New("no data inserted")
	}
	return templates, nil
}

func GetSMSTemplateByVirtualID(ctx context.Context, virtualId uint) (*SMSTemplate, error) {
	template := new(SMSTemplate)
	result := OpenDB.WithContext(ctx).Table(template.TableName()).Where("virtual_id = ?", virtualId).Scan(template)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, gorm.ErrRecordNotFound
	}
	return template, nil
}

func GetAllTemplate(ctx context.Context) (*[]SMSTemplate, error) {
	templates := new([]SMSTemplate)
	result := OpenDB.Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}
	return templates, nil
}
