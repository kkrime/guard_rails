package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RepositoryController interface {
	AddRepository(ctx *gin.Context)
	GetRepository(ctx *gin.Context)
	UpdateRepository(ctx *gin.Context)
	DeleteRepository(ctx *gin.Context)
}

type ScanController interface {
	Scan(ctx *gin.Context)
	GetScans(ctx *gin.Context)
}

type getLogger interface {
	getLogger() *logrus.Logger
}
