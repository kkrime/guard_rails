package server

import (
	"guard_rails/config"
	"guard_rails/controller"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	requestid "github.com/sumit-tembe/gin-requestid"
)

func Init(database *sqlx.DB, config *config.Config, log *logrus.Logger) (*gin.Engine, error) {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(requestid.RequestID(nil))
	router.Use(gin.LoggerWithConfig(requestid.GetLoggerConfig(nil, nil, nil)))

	v1 := router.Group("v1")

	// repository
	repositoryGroup := v1.Group("repository")
	repositoryController := controller.NewRepoitoryController(database, log)

	// add repository
	repositoryGroup.POST("", repositoryController.AddRepository)
	// get repository
	repositoryGroup.GET("/:name", repositoryController.GetRepository)
	// update repository
	repositoryGroup.PUT("", repositoryController.UpdateRepository)
	// delete repository
	repositoryGroup.DELETE("/:name", repositoryController.DeleteRepository)

	// scan
	scanGroup := repositoryGroup.Group("scan")
	scanController, err := controller.NewScanController(database, config, log)
	if err != nil {
		return nil, err
	}

	scanGroup.POST("/:name", scanController.QueueScan)
	scanGroup.GET("/:name", scanController.GetScans)

	return router, nil
}
