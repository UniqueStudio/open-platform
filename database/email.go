package database

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type EmailTemplate struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ID          uint   `gorm:"column:id;primarykey" json:"template_id"`
	TemplateID  string `gorm:"column:template_id" json:"-"`
	Content     string `gorm:"column:content" json:"content"`
}

type EmailHTMLTemplate struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ID          uint   `gorm:"column:id;primarykey" json:"template_id"`
	TemplateID  string `gorm:"column:template_id" json:"-"`
	Style 		string `gorm:"column:style" json:"-"`
	Content     string `gorm:"column:content" json:"content"`
}

func (t *EmailTemplate) TableName() string {
	return "email_template"
}

func (t *EmailHTMLTemplate) TableName() string {
	return "email_HTML_template"
}

func GetEmailTemplateByID(ctx context.Context, id uint) *EmailTemplate {
	template := new(EmailTemplate)
	OpenDB.WithContext(ctx).Table(template.TableName()).Where("id = ?", id).Scan(template)
	return template
}

func GetAllEmailTemplate(ctx context.Context) (*[]EmailTemplate, error) {
	templates := new([]EmailTemplate)

	if err := OpenDB.WithContext(ctx).Find(templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}
