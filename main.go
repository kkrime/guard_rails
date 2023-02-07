package main

import (
	"fmt"
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
	fmt.Println(config)
	if err != nil {
		return err
	}

	db, err := db.Init(&config.Postgres)
	if err != nil {
		return err
	}

	r := server.Init(db, log)
	log.Info("starting server")

	return r.Run()
}
