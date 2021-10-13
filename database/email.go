package database

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type EmailTemplate struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	TemplateID   uint32 `gorm:"column:template_id" json:"template_id"`
	TemplateName string `gorm:"column:template_id" json:"-"`
	Content      string `gorm:"column:content" json:"content"`
	ParamNumber  uint32 `gorm:"column:param_number" json:"param_number"`
}

func (t *EmailTemplate) TableName() string {
	return "email_template"
}

func GetEmailTemplateByID(ctx context.Context, id uint) *EmailTemplate {
	template := new(EmailTemplate)
	OpenDB.WithContext(ctx).Table(template.TableName()).Where("template_id = ?", id).Scan(template)
	return template
}

func GetAllEmailTemplate(ctx context.Context) (*[]EmailTemplate, error) {
	templates := new([]EmailTemplate)
	if err := OpenDB.WithContext(ctx).Find(templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func InsertEmailTemplate(ctx context.Context, templates *[]EmailTemplate) (*[]EmailTemplate, error) {
	result := OpenDB.Create(templates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, errors.New("no data inserted")
	}
	return templates, nil
}
