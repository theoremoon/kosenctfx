package config

import (
	"fmt"
	"os"
)

type Config struct {
	Dbdsn                string
	Addr                 string
	RedisAddr            string
	Front                string
	Email                string
	MailServer           string
	MailPassword         string
	AdminWebhookURL      string
	SolveCheckWebhookURL string
	SystemWebhookURL     string
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

	redisAddr, err := getEnv("REDIS")
	if err != nil {
		redisAddr = ":6379"
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

	adminWebhookURL, _ := getEnv("ADMIN_WEBHOOK")
	solveCheckWebhookURL, _ := getEnv("SOLVE_WEBHOOK")
	systemWebhookURL, _ := getEnv("SYSTEM_WEBHOOK")

	return &Config{
		Dbdsn:                dbdsn,
		Addr:                 addr,
		RedisAddr:            redisAddr,
		Front:                front,
		Email:                mailaccount,
		MailServer:           mailserver,
		MailPassword:         mailpassword,
		AdminWebhookURL:      adminWebhookURL,
		SolveCheckWebhookURL: solveCheckWebhookURL,
		SystemWebhookURL:     systemWebhookURL,
	}, nil
}
