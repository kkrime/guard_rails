package controller

import (
	"guard_rails/model"
	"guard_rails/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type repositoryController struct {
	repositoryServiceProvider service.RepositoryServiceProvider
	log                       *logrus.Logger
}

func (rc *repositoryController) getLogger() *logrus.Logger {
	return rc.log
}

func NewRepoitoryController(db *sqlx.DB, log *logrus.Logger) RepositoryController {

	repositoryServiceProvider := service.NewRepositoryServiceProvider(db)
	return &repositoryController{
		repositoryServiceProvider: repositoryServiceProvider,
		log:                       log,
	}
}

func (rc *repositoryController) AddRepository(ctx *gin.Context) {
	repository := &model.Repository{}

	err := ctx.BindJSON(repository)
	if err != nil {
		AbortAndError(ctx, err)
		return
	}

	log := createLogger(rc, ctx)
	repositoryServiceInstance := rc.repositoryServiceProvider.NewRepositoryServiceInstance(log)

	err = repositoryServiceInstance.AddRepository(ctx, repository)
	if err != nil {
		AbortAndError(ctx, err)
		return
	}

	ctx.JSON(201, getSuccessResponse())
}

func (rc *repositoryController) GetRepository(ctx *gin.Context) {

	repositoryName := ctx.Param("name")
	repositoryName = strings.ToLower(repositoryName)

	log := createLogger(rc, ctx)
	repositoryServiceInstance := rc.repositoryServiceProvider.NewRepositoryServiceInstance(log)

	repository, err := repositoryServiceInstance.GetRepository(ctx, repositoryName)
	if err != nil {
		AbortAndError(ctx, err)
		return
	}

	response := getSuccessResponse()
	response.Data = repository

	ctx.JSON(200, response)
}

func (rc *repositoryController) UpdateRepository(ctx *gin.Context) {
	repository := &model.Repository{}

	err := ctx.BindJSON(repository)
	if err != nil {
		AbortAndError(ctx, err)
		return
	}
	repository.Name = strings.ToLower(repository.Name)

	log := createLogger(rc, ctx)
	repositoryServiceInstance := rc.repositoryServiceProvider.NewRepositoryServiceInstance(log)

	err = repositoryServiceInstance.UpdateRepository(ctx, repository)
	if err != nil {
		AbortAndError(ctx, err)
		return
	}

	ctx.JSON(200, getSuccessResponse())
}

func (rc *repositoryController) DeleteRepository(ctx *gin.Context) {

	repositoryName := ctx.Param("name")
	repositoryName = strings.ToLower(repositoryName)

	log := createLogger(rc, ctx)
	repositoryServiceInstance := rc.repositoryServiceProvider.NewRepositoryServiceInstance(log)

	err := repositoryServiceInstance.DeleteRepository(ctx, repositoryName)
	if err != nil {
		AbortAndError(ctx, err)
		return
	}

	ctx.JSON(200, getSuccessResponse())
}
