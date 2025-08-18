package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(env string) (string, error) {
	if _, ok := os.LookupEnv(env); !ok {
		if err := godotenv.Load(); err != nil {
			return "", err
		}
	}

	return os.Getenv(env), nil
}
