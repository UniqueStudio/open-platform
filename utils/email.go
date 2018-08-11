package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// RenderHTML is a func to gen html from html template
func RenderHTML(name, content string) (html string) {
	t, err := template.ParseFiles("./static/html/email.html")
	if err != nil {
		fmt.Println(err)
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

	return body.String()
}
