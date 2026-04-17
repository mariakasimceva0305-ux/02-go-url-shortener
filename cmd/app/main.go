package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"go-url-shortener/internal/api"
	"go-url-shortener/internal/repository"
	"go-url-shortener/internal/service"
	"go-url-shortener/internal/worker"
)

func main() {
	dbPath := envOrDefault("DB_PATH", "data.db")
	baseURL := envOrDefault("BASE_URL", "http://localhost:8080")
	port := envOrDefault("PORT", "8080")

	log.Printf("starting storage for %s", dbPath)
	linkRepo := repository.NewMemoryLinkRepository()
	eventRepo := repository.NewMemoryEventRepository()

	linkService := service.NewLinkService(linkRepo, eventRepo, baseURL)
	worker.StartCleanup(linkService, 1*time.Minute)

	router := gin.Default()
	handler := api.NewHandler(linkService)
	handler.Register(router)

	log.Printf("starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
