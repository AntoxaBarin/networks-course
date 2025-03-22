package main

import (
	"flag"
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	text = "This is the email body."
	html = `<!DOCTYPE html>
<html>
<head>
<title>Моя первая веб-страница</title>
</head>
<body>
<h1>Добро пожаловать на мою первую веб-страницу!</h1>
<p>Это пример простого HTML-документа.</p>
<ul>
  <li>1. C++</li>
  <li>2. Ada</li>
  <li>3. Алгол 68</li>
</ul>
</body>
</html>`
)

var (
	from     string
	to       string
	password string
	msg_type string
)

func main() {
	flag.StringVar(&from, "from", "dummy_mail@aboba.com", "Sender email")
	flag.StringVar(&to, "to", from, "Recepient email")
	flag.StringVar(&password, "pass", "strong_password", "Sender mail password")
	flag.StringVar(&msg_type, "type", "text", "Type of message: text or html")
	flag.Parse()

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello!")

	if msg_type == "text" {
		m.SetBody("text/plain", text)
	} else if msg_type == "html" {
		m.SetBody("text/html", html)
	} else {
		panic("Bad message type!")
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully")
}
