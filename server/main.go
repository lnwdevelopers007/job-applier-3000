package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	router := controller.NewRouter()
	router.Run(os.Getenv("PORT"))
}
