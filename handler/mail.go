package handler

import (
	"fmt"
	"net/http"
	"net/mail"
	"open-platform/utils"

	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/gin-gonic/gin"
)

type sender struct {
	To      string `json:"to"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// SendMailHandler is a func to handle send email template requests
func SendMailHandler(c *gin.Context) {
	var data sender
	c.BindJSON(&data)

	to := data.To
	name := data.Name
	content := data.Content
	subject := data.Subject

	if name == "" || to == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Post data missing parameter"})
		return
	}

	fmt.Println("send email")

	err := SendToMail(utils.AppConfig.SMTP.Sender, utils.AppConfig.SMTP.Password, utils.AppConfig.SMTP.Host, to, subject, utils.RenderHTML(name, content), "html")

	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
		c.JSON(http.StatusConflict, gin.H{"msg": err.Error(), "code": http.StatusBadRequest})
	} else {
		fmt.Println("Send mail success!")
		c.JSON(http.StatusOK, gin.H{"msg": "OK", "code": http.StatusOK})
	}
}

// SendToMail is a function to handle send email smtp requests
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	auth := sasl.NewPlainClient("", user, password)
	fromName := "联创团队"
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain; charset=UTF-8"
	}
	msg := strings.NewReader("To: " + to + "\r\nReply-To: " + "contact@hustunique.com" + "\r\nFrom: " + fromName + " <" + user + ">\r\nSubject: " + encodeRFC2047(subject) + "\r\n" + contentType + "\r\n\r\n" + body)

	sendTo := strings.Split(to, ";")

	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), "<@>")
}
