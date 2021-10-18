package utils

import (
	"context"
	"fmt"
	"github.com/UniqueStudio/open-platform/config"
	"github.com/UniqueStudio/open-platform/database"
	"github.com/UniqueStudio/open-platform/pkg"
	gomail "github.com/go-mail/gomail"
	"log"
)

type MailClient struct {
	Host        string
	Port        int
	Sender      string
	Password    string
	Name        string
	ContentType string
}

var mailClient *MailClient

func SetupEmailClient() error {
	mailClient = &MailClient{
		Host:        config.Config.Email.SMTP.Host,
		Port:        config.Config.Email.SMTP.Port,
		Sender:      config.Config.Email.SMTP.Sender,
		Password:    config.Config.Email.SMTP.Password,
		Name:        config.Config.Email.SMTP.Name,
		ContentType: "text/html; charset=UTF-8",
	}
	return nil
}

func SendSingleEmail(ctx context.Context, to string, templateID uint, subject string, paramSet []string) (*pkg.EmailSendStatus, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", mailClient.Name)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	template := database.GetEmailTemplateByID(ctx, templateID)
	var set []interface{}
	for _, v := range paramSet {
		set = append(set, v)
	}
	m.SetBody(mailClient.ContentType, fmt.Sprintf(template.Content, set...))

	d := gomail.NewDialer(mailClient.Host, mailClient.Port, mailClient.Sender, mailClient.Password)

	resp := &pkg.EmailSendStatus{
		To:      to,
		Content: fmt.Sprintf(template.Content, set...),
	}

	err := d.DialAndSend(m)
	if err != nil {
		resp.Err = err.Error()
		return resp, err
	}
	m.Reset()
	return resp, nil
}

func SendMultipleEmail(ctx context.Context, to []string, templateID uint, subject string, paramSet [][]string) ([]*pkg.EmailSendStatus, error) {
	status := make([]*pkg.EmailSendStatus, 0)
	for i := 0; i < len(to); i++ {
		resp, err := SendSingleEmail(ctx, to[i], templateID, subject, paramSet[i])
		if err != nil {
			log.Println(err)
		}
		status = append(status, resp)
	}
	return status, nil
}
