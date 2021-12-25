package config

import (
	"fmt"
	"os"
)

type Config struct {
	Dbdsn              string
	Addr               string
	RedisAddr          string
	Front              string
	MailFake           bool
	Email              string
	MailServer         string
	MailPassword       string
	AdminWebhookURL    string
	SolveLogWebhookURL string
	TaskOpenWebhookURL string
	AdminToken         string
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
	dbdsn += "?parseTime=true&charset=utf8mb4&collation=utf8mb4_bin"

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
	mailFake := false
	if _, err := getEnv("MAIL_FAKE"); err == nil {
		mailFake = true
	}

	adminWebhookURL, _ := getEnv("ADMIN_WEBHOOK")
	solveLogWebhookURL, _ := getEnv("SOLVE_WEBHOOK")
	taskOpenWebhookURL, _ := getEnv("TASK_OPEN_WEBHOOK")
	adminToken, _ := getEnv("ADMIN_TOKEN")

	return &Config{
		Dbdsn:              dbdsn,
		Addr:               addr,
		RedisAddr:          redisAddr,
		Front:              front,
		MailFake:           mailFake,
		Email:              mailaccount,
		MailServer:         mailserver,
		MailPassword:       mailpassword,
		AdminWebhookURL:    adminWebhookURL,
		SolveLogWebhookURL: solveLogWebhookURL,
		TaskOpenWebhookURL: taskOpenWebhookURL,
		AdminToken:         adminToken,
	}, nil
}
