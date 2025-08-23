package controller

import "github.com/gin-gonic/gin"

// IController is the interface of all controllers.
type IController[T any] interface {
	Create() gin.HandlerFunc
	RetrieveAll() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

// Query() method to be added.
