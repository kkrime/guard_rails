package controller

import (
	"guard_rails/config"
	"guard_rails/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type scanController struct {
	scanServiceProvider service.ScanServiceProvider
	log                 *logrus.Logger
}

func NewScanController(db *sqlx.DB, config *config.Config, log *logrus.Logger) (ScanController, error) {

	ScanServiceProvider, err := service.NewScanServiceProvider(db, config)
	if err != nil {
		return nil, err
	}
	return &scanController{
		scanServiceProvider: ScanServiceProvider,
		log:                 log,
	}, nil
}

func (sc *scanController) getLogger() *logrus.Logger {
	return sc.log
}

func (sc *scanController) QueueScan(ctx *gin.Context) {
	repositoryName := ctx.Param("name")
	repositoryName = strings.ToLower(repositoryName)

	log := createLogger(sc, ctx)
	scanServiceInstance := sc.scanServiceProvider.NewScanServiceInstance(log)

	err := scanServiceInstance.QueueScan(ctx, repositoryName)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := getSuccessResponse()

	ctx.JSON(200, response)
}

func (sc *scanController) GetScans(ctx *gin.Context) {
	repositoryName := ctx.Param("name")
	repositoryName = strings.ToLower(repositoryName)

	log := createLogger(sc, ctx)
	scanServiceInstance := sc.scanServiceProvider.NewScanServiceInstance(log)

	scans, err := scanServiceInstance.GetScan(ctx, repositoryName)
	if err != nil {
		ctx.Error(err)
		return
	}

	output := transsformScans(scans)
	response := getSuccessResponse()
	response.Data = output

	ctx.JSON(200, response)
}
