package server

import (
	"encoding/json"
	"guard_rails/controller"
	"guard_rails/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

const (
	Invalid_Input  = "invalid input data"
	Internal_Error = "internal error"
)

func Init(db *sqlx.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(ErrorHandler)

	v1 := router.Group("v1")
	// repository
	repositoryGroup := v1.Group("repository")
	repositoryController := controller.NewRepoitoryController(db)

	// add repository
	repositoryGroup.POST("", repositoryController.AddRepository)
	// get repository
	repositoryGroup.GET("/:name", repositoryController.GetRepository)
	// update repository
	repositoryGroup.PUT("", repositoryController.UpdateRepository)
	// delete repository
	repositoryGroup.DELETE("/:name", repositoryController.DeleteRepository)

	return router
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		switch err := err.Err.(type) {
		case validator.ValidationErrors:
			c.JSON(400, gin.H{"error": controller.GetFailResponse(Invalid_Input)})
		case *json.InvalidUnmarshalError:
			c.JSON(400, gin.H{"error": controller.GetFailResponse(Invalid_Input)})
		case *errors.RestError:
			c.JSON(err.Code, controller.GetFailResponse(err.Err))
		default:
			c.JSON(500, controller.GetFailResponse(Internal_Error))
		}
	}
}
