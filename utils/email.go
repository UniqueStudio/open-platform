package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// RenderHTML is a func to generate html from html template
func RenderHTML(name, content string) (html string, err error) {
	t, err := template.ParseFiles("./static/html/email.html")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var body bytes.Buffer
	t.Execute(&body, struct {
		Name    string
		Content string
		Year    int
	}{
		Name:    name,
		Content: content,
		Year:    time.Now().Year(),
	})

	//// Configure hermes by setting a theme and your product info
	//h := hermes.Hermes{
	//	// Optional Theme
	//	// Theme: new(Default)
	//	Product: hermes.Product{
	//		// Appears in header & footer of e-mails
	//		Name: "Unique Studio",
	//		Link: "https://www.hustunique.com/",
	//		// Optional product logo
	//		Logo: "https://storage.fredliang.cn/web/studio/Logo.png",
	//	},
	//}
	//
	//email := hermes.Email{
	//	Body: hermes.Body{
	//		Name:         name,
	//		FreeMarkdown: hermes.Markdown(content),
	//	},
	//}
	//
	//emailBody, err := h.GenerateHTML(email)
	//if err != nil {
	//	return "", err // Tip: Handle error with something else than a panic ;)
	//}

	return body.String(), nil
}
