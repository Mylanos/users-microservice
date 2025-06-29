package main

import (
	"log"
	"users-backend/pkg/api"
	"users-backend/pkg/config"
	"users-backend/pkg/services"
	"users-backend/pkg/storage"

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
	api := api.NewAPIServer(":8080", service)
	if err := api.Router(); err != nil {
		log.Fatalf("FATAL: could not start server: %v", err)
	}
}
