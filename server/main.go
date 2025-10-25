// @title Job Applier 3000 API
// @version 1.0
// @description REST API for Job Applier 3000 platform
// @BasePath /
// @schemes http https
package main

import (
	"os"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
)

func main() {
	router := controller.NewRouter()
	router.Run(os.Getenv("SERVER_ADDR"))
}
