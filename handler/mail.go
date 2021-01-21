package handler

import (
	"fmt"
	"net/http"
	"net/mail"
	"open-platform/utils"

	"strings"

	sasl "github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/gin-gonic/gin"
)

type sender struct {
	BccTo	[]string `json:"bccto"`
	CcTo	[]string `json:"ccto"`
	To      string   `json:"to"`
	ToList  []string `json:"toList"`
	Name    string   `json:"name"`
	Subject string   `json:"subject"`
	Content string   `json:"content"`
}

// SendMailHandler is a func to handle send email template requests
func SendMailHandler(c *gin.Context) {
	var data sender
	c.BindJSON(&data)

	bccto := data.BccTo
	ccto := data.CcTo
	to := data.To
	toList := data.ToList
	name := data.Name
	content := data.Content
	subject := data.Subject

	if ((to == "") && (len(toList) == 0)) || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Post data missing parameter"})
		return
	}

	renderContent, err := utils.RenderHTML(name, content)

	if err != nil {
		fmt.Println("Send mail error1!")
		fmt.Println(err)
		c.JSON(http.StatusConflict, gin.H{"message": err.Error(), "code": http.StatusBadRequest})
	}

	err = SendToMail(
		utils.AppConfig.SMTP.Sender,
		utils.AppConfig.SMTP.Password,
		utils.AppConfig.SMTP.Host,
		subject, renderContent, "html", utils.RemoveDuplicate(append(toList, to)), ccto, bccto)

	if err != nil {
		fmt.Println("Send mail error2!")
		fmt.Println(err)
		c.JSON(http.StatusConflict, gin.H{"message": err.Error(), "code": http.StatusBadRequest})
	} else {
		fmt.Println("Send mail success!")
		c.JSON(http.StatusOK, gin.H{"message": "OK", "code": http.StatusOK})
	}
}

// SendToMail is a function to handle send email smtp requests
func SendToMail(user, password, host, subject, body, mailtype string, to, ccto, bccto []string) error {
	auth := sasl.NewPlainClient("", user, password)
	fromName := "联创团队"
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain; charset=UTF-8"
	}

	toAdress := strings.Join(to, ",")
	cctoAdress := strings.Join(ccto, ",")
	bcctoAdress := strings.Join(bccto,",")
	msg := strings.NewReader("To: " + toAdress + "\r\nReply-To: " + "contact@hustunique.com" +  "\r\nCc: " + cctoAdress + "\r\nBcc: " + bcctoAdress +  "\r\nFrom: " + fromName + " <" + user + ">\r\nSubject: " + encodeRFC2047(subject) + "\r\n" + contentType + "\r\n\r\n" + body)

	to = MergeSlice(to, ccto)
	to = MergeSlice(to, bccto)
	err := smtp.SendMail(host, auth, user, to, msg)
	return err
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Name: String, Address: ""}
	return strings.Trim(addr.String(), "<@>")
}

func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}