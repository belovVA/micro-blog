package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const pathEnv = ".env"
const cfgPathDefault = "./configs/config.yaml"

type HTTPConfig interface {
	GetPort() string
	GetHost() string
	GetTimeout() time.Duration
	GetIdleTimeout() time.Duration
}

func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}

func LoadConfig() (string, error) {
	if err := LoadEnv(pathEnv); err != nil {
		return "", err
	}

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = cfgPathDefault
	}

	if _, err := os.Stat(cfgPath); err != nil {
		return "", fmt.Errorf("%s file not found", cfgPath)
	}

	return cfgPath, nil
}
