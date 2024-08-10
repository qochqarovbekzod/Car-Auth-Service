package main

import (
	"auth/api"
	"auth/api/handler"
	"auth/cmd/server"
	"auth/config"
	"auth/logs"
	"auth/service"
	"auth/storage"
	"auth/storage/postgres"
	"log"
	"sync"
)

func main() {
	logs.InitLogger()

	log.Println("Starting the server ...")

	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
		log.Fatal(err)
		return
	}

	defer db.Close()

	cfg := config.Load()
	handler := handler.NewHandler(service.NewService(storage.NewStorage(db,logs.Logger), logs.Logger))
	router := api.NewRouter(handler)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		logs.Logger.Info("server is running", "port:", cfg.HTTP_PORT)
		err := router.Run(cfg.HTTP_PORT)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
			log.Fatal(err)
		}
	}()

	server.RunServer(storage.NewStorage(db, logs.Logger))
	wg.Wait()

}
