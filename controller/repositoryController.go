package controller

import (
	"fmt"
	"guard_rails/model"
	"guard_rails/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type repositoryController struct {
	repositoryService service.RepositoryService
}

func NewRepoitoryController(db *sqlx.DB) RepositoryController {

	repositoryService := service.NewRepositoryService(db)
	return &repositoryController{
		repositoryService: repositoryService,
	}
}

func (rc *repositoryController) AddRepository(ctx *gin.Context) {
	repository := &model.Repository{}

	err := ctx.BindJSON(repository)
	if err != nil {
		return
	}

	err = rc.repositoryService.AddRepository(ctx, repository)
	if err != nil {
		ctx.AbortWithError(-1, err)
		return
	}

	ctx.JSON(201, getSuccessResponse())
}

func (rc *repositoryController) GetRepository(ctx *gin.Context) {

	repositoryName := ctx.Param("name")
	repositoryName = strings.ToLower(repositoryName)

	repository, err := rc.repositoryService.GetRepository(ctx, repositoryName)
	if err != nil {
		ctx.AbortWithError(-1, err)
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
		return
	}
	repository.Name = strings.ToLower(repository.Name)

	err = rc.repositoryService.UpdateRepository(ctx, repository)
	if err != nil {
		ctx.AbortWithError(-1, err)
		return
	}

	ctx.JSON(200, getSuccessResponse())
}

func (rc *repositoryController) DeleteRepository(ctx *gin.Context) {

	repositoryName := ctx.Param("name")
	repositoryName = strings.ToLower(repositoryName)

	err := rc.repositoryService.DeleteRepository(ctx, repositoryName)
	if err != nil {
		fmt.Printf("err = %+v\n", err)
		ctx.AbortWithError(-1, err)
		return
	}

	ctx.JSON(200, getSuccessResponse())
}
