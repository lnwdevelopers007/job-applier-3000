package controller

import "github.com/gin-gonic/gin"

// IController is the interface of all controllers.
type IController[T any] interface {
	Create(*gin.Context)
	RetrieveAll(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

// Query() method to be added.
