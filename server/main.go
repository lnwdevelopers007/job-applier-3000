package main

import (
	"os"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
)

func main() {
	router := controller.NewRouter()
	router.Run(os.Getenv("SERVER_ADDR"))
}
