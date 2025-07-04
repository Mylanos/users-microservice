package main

import (
	"log"
	"users-microservice/pkg/api"
	"users-microservice/pkg/config"
	"users-microservice/pkg/services"
	"users-microservice/pkg/storage"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	storageImpl, err := storage.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatalf("FATAL: failed to create a storage: %s", err)
	}
	service, err := services.NewUserService(storageImpl)
	if err != nil {
		log.Fatalf("FATAL: failed to create a UserService: %s", err)
	}
	apiServer := api.NewAPIServer(":8080", service, cfg)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("FATAL: could not start server: %v", err)
	}
}
