package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hanlinthedev/file-service/internal/config"
	"github.com/hanlinthedev/file-service/internal/database"
	"github.com/hanlinthedev/file-service/internal/handlers"
	"github.com/hanlinthedev/file-service/internal/models"
	"github.com/hanlinthedev/file-service/internal/repositories"
	"github.com/hanlinthedev/file-service/internal/services"
	"github.com/hanlinthedev/file-service/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No ENV Found, %s", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.File{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.MaxMultipartMemory = cfg.MaxFileSize

	repo := repositories.NewPgFileRepository(db)
	store := storage.NewLocalStorage(cfg.StoragePath)
	service := services.NewFileService(repo, store, cfg.MaxFileSize)
	handler := handlers.NewFileHandler(service)

	r.GET("/health", handlers.HealthCheck)
	r.POST("/files/upload", handler.Upload)
	r.DELETE("/files/delete/:id", handler.Delete)
	r.GET("/files/:id", handler.GetById)
	r.GET("/files/:id/download", handler.Download)

	log.Printf("Listening on port %s", cfg.Port)
	if err = r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
