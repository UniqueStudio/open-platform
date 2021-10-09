package utils

import (
	"github.com/UniqueStudio/open-platform/pkg"
	"net/smtp"
	"crypto/tls"
)


type EmailClient struct {
	host string
	port uint
	secret string
}

var (
	client *EmailClient
)

func SetupEmailClient() error {

}


func CheckConn() error {
	
}

func SendSingleEmail(to string, paramSet []string, templateID uint) (*pkg.EmailSendStatus, error) {
	
}

func SendMultipleEmail(to []string, paramSet []string, templateID uint) ([]*pkg.EmailSendStatus, error) {

}

func GetEmailTemplates()



func (client *EmailClient) Close() error {

}

