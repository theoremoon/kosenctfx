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
	BucketEndpoint     string
	BucketRegion       string
	BucketAccessKey    string
	BucketSecretKey    string
	BucketName         string
	InsecureBucket     bool
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

	bucketEndpoint, _ := getEnv("BUCKET_ENDPOINT")
	bucketRegion, _ := getEnv("BUCKET_REGION")
	bucketAccessKey, _ := getEnv("BUCKET_ACCESS_KEY")
	bucketSecretKey, _ := getEnv("BUCKET_SECRET_KEY")
	bucketName, _ := getEnv("BUCKET_NAME")
	insecureBucket := false
	if _, err := getEnv("BUCKET_INSECURE"); err == nil {
		insecureBucket = true
	}
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
		BucketEndpoint:     bucketEndpoint,
		BucketRegion:       bucketRegion,
		BucketAccessKey:    bucketAccessKey,
		BucketSecretKey:    bucketSecretKey,
		BucketName:         bucketName,
		InsecureBucket:     insecureBucket,
		AdminToken:         adminToken,
	}, nil
}
