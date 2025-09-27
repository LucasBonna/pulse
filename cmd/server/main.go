package main

import (
	"log"

	"lucasbonna/pulse/internal/api"
	"lucasbonna/pulse/internal/storage"
)

func main() {
	dbInstance, err := storage.NewSQLiteDB()
	if err != nil {
		log.Fatal("error creating db")
	}

	httpServer := api.NewServer(dbInstance)

	err = httpServer.Start(":3333")
	if err != nil {
		log.Fatal("error starting http server at port 3333", err)
	}

}
