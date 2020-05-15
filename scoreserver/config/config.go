package config

import (
	"fmt"
	"os"
)

type Config struct {
	Dbdsn string
	Addr  string
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

	addr, err := getEnv("ADDR")
	if err != nil {
		addr = ":5000"
	}
	return &Config{
		Dbdsn: dbdsn,
		Addr:  addr,
	}, nil
}
