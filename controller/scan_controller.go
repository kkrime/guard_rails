package controller

import (
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

func NewScanController(db *sqlx.DB, log *logrus.Logger) ScanController {

	ScanServiceProvider := service.NewScanServiceProvider(db)
	return &scanController{
		scanServiceProvider: ScanServiceProvider,
		log:                 log,
	}
}

func (sc *scanController) getLogger() *logrus.Logger {
	return sc.log
}

func (sc *scanController) Scan(ctx *gin.Context) {
	repositoryName := ctx.Param("name")
	repositoryName = strings.ToLower(repositoryName)

	log := createLogger(sc, ctx)
	scanServiceInstance := sc.scanServiceProvider.NewScanServiceInstance(log)

	err := scanServiceInstance.Scan(ctx, repositoryName)
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
