package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(env string) string {
	if _, ok := os.LookupEnv(env); !ok {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
			return ""
		}
	}

	return os.Getenv(env)
}

func LoadOAuthRedirectURI(protocol string, provider string) string {
	return protocol + "://" + LoadEnv("SERVER_ADDR") + "/auth/" + provider + "/callback"
}
