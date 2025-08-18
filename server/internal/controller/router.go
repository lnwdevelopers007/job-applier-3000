package controller

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/jobs", GetJobs)
	router.POST("/jobs", PostJobs)
	return router
}
