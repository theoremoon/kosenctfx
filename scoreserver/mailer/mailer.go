package mailer

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

type Mailer interface {
	Send(to, subject, body string) error
}

type mailer struct {
	account  string
	password string
	server   string
	host     string
}

func New(server, account, password string) (Mailer, error) {
	host := strings.Split(server, ":")
	if len(host) != 2 {
		return nil, errors.New("server address format: <HOST>:<PORT>")
	}

	return &mailer{
		account:  account,
		password: password,
		server:   server,
		host:     host[0],
	}, nil
}

func (m *mailer) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.account, m.password, m.host)
	msg := "From: " + m.account + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	err := smtp.SendMail(m.server, auth, m.account, []string{to}, []byte(msg))
	if err != nil {
		fmt.Errorf("%w", err)
	}

	return nil
}
