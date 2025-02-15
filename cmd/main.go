package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/database"
	"github.com/RazuOff/MerchShop/internal/handler"
	"github.com/RazuOff/MerchShop/internal/repository"
	"github.com/RazuOff/MerchShop/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")
	config := config.NewConfig()

	db, err := database.NewPostgre(config)
	if err != nil {
		log.Fatalf("Data base connection failed\nError: %s", err.Error())
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service, config)

	r := handler.InitRoutes()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()
	log.Println("Server started on :8080")

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exited properly")
}
