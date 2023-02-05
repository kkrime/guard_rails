package main

import (
	"guard_rails/db"
	"guard_rails/server"

	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	db, err := db.Init()

	if err != nil {
		return err
	}
	defer db.Close()

	r := server.Init(db)
	return r.Run()
}
