// @title Job Applier 3000 API
// @version 1.0
// @description REST API for Job Applier 3000 platform
// @BasePath /
// @schemes http https
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
)

func main() {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	slog.SetLogLoggerLevel(slog.LevelInfo)

	router := controller.NewRouter()
	router.Run(os.Getenv("SERVER_ADDR"))
}
