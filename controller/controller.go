package controller

import "github.com/gin-gonic/gin"

type RepositoryController interface {
	AddRepository(ctx *gin.Context)
	GetRepository(ctx *gin.Context)
	UpdateRepository(ctx *gin.Context)
	DeleteRepository(ctx *gin.Context)
}
