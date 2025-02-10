package main

import (
	"log"

	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/database"
	"github.com/RazuOff/MerchShop/internal/handler"
	"github.com/RazuOff/MerchShop/internal/repository"
	"github.com/RazuOff/MerchShop/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("No .env file found")
	}
	config := config.NewConfig()

	db, err := database.NewPostgre(config)
	if err != nil {
		log.Fatalf("Data base connection failed\nError: %s", err.Error())
	}

	repository := repository.NewRepository(db)

	service := service.NewService(repository)
	handler := handler.NewHandler(service, config)

	r := handler.InitRoutes()

	r.Run()

}
