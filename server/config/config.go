package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv(env string) string {
	if _, ok := os.LookupEnv(env); !ok {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
			return "false"
		}
	}

	return os.Getenv(env)
}

func LoadCallbackURI(protocol string, provider string) string {
	return LoadEnv("OAUTH_REDIRECT_URL")
}

func LoadBoolean(env string) bool {
	res, err := strconv.ParseBool(LoadEnv(env))
	if err != nil {
		log.Fatal(err)
		return false
	}
	return res
}

func LoadInt(env string) int {
	res, err := strconv.ParseInt(LoadEnv(env), 10, 32)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return int(res)
}
