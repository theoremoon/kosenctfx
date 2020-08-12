package config

import (
	"fmt"
	"os"
)

type Config struct {
	Dbdsn        string
	Addr         string
	Front        string
	Email        string
	MailServer   string
	MailPassword string
}

func getEnv(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("Environmental vairable '%s' is required", name)
	}
	return val, nil
}

func Load() (*Config, error) {
	dbdsn, err := getEnv("DBDSN")
	if err != nil {
		return nil, err
	}
	dbdsn += "?parseTime=true"

	addr, err := getEnv("ADDR")
	if err != nil {
		addr = ":5000"
	}

	front, err := getEnv("FRONT")
	if err != nil {
		return nil, err
	}

	mailaccount, err := getEnv("MAIL")
	if err != nil {
		return nil, err
	}

	mailserver, err := getEnv("MAIL_SERVER")
	if err != nil {
		return nil, err
	}

	mailpassword, err := getEnv("MAIL_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &Config{
		Dbdsn:        dbdsn,
		Addr:         addr,
		Front:        front,
		Email:        mailaccount,
		MailServer:   mailserver,
		MailPassword: mailpassword,
	}, nil
}
