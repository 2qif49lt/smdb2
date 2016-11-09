package main

import (
	"encoding/base64"
	"fmt"
	//	"github.com/scorredoira/email"
	"net/smtp"
	"strings"
)

func checkAdminMail(email string) bool {
	for _, v := range cfg.Admin {
		if v == email {
			return true
		}
	}
	return false
}

/*
func sendMail(to, title, content string) error {
	m := email.NewHTMLMessage(title, content)
	m.From = mail.Address{Name: "smdb2", Address: cfg.Smtp.Email}
	m.To = []string{to}

	auth := smtp.PlainAuth("", cfg.Smtp.Email, cfg.Smtp.Pwd, cfg.Smtp.Srv)
	return email.Send(fmt.Sprintf(`%s:%d`, cfg.Smtp.Srv, cfg.Smtp.Port), auth, m)
}
*/

func sendMail2(from, to, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.Smtp.Email, cfg.Smtp.Pwd, cfg.Smtp.Srv)
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	msg := []byte("To: " + to + "\r\nFrom: " + from +
		"<" + cfg.Smtp.Email + ">\r\nSubject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8 \r\nContent-Transfer-Encoding: base64\r\n\r\n" + b64.EncodeToString([]byte(body)))
	send_to := strings.Split(to, "|")
	err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Smtp.Srv, cfg.Smtp.Port),
		auth,
		cfg.Smtp.Email,
		send_to,
		msg)
	return err
}
