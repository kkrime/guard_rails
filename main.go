package main

import (
	"guard_rails/config"
	"guard_rails/db"
	"guard_rails/logger"
	"guard_rails/server"

	"github.com/sirupsen/logrus"
)

func main() {

	log := logger.CreateNewLogger()
	log.ReportCaller = true

	if err := run(log); err != nil {
		log.Fatal(err)
	}
}

func run(log *logrus.Logger) error {

	config, err := config.ReadConfig()
	if err != nil {
		return err
	}

	database, err := db.Init(&config.Postgres)
	if err != nil {
		return err
	}
	log.Info("connected to database")

	r, err := server.Init(database, config, log)
	if err != nil {
		return err
	}
	log.Info("starting server")

	return r.Run()
}
