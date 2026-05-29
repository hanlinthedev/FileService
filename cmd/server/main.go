package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hanlinthedev/file-service/internal/config"
	"github.com/hanlinthedev/file-service/internal/database"
	"github.com/hanlinthedev/file-service/internal/handlers"
	"github.com/hanlinthedev/file-service/internal/logger"
	"github.com/hanlinthedev/file-service/internal/middleware"
	"github.com/hanlinthedev/file-service/internal/models"
	"github.com/hanlinthedev/file-service/internal/repositories"
	"github.com/hanlinthedev/file-service/internal/services"
	"github.com/hanlinthedev/file-service/internal/storage"
	"github.com/hanlinthedev/file-service/internal/workers"
	"github.com/joho/godotenv"
)

func main() {
	log := logger.AppLogger()

	allowedMimeTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"application/pdf": true,
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".pdf":  true,
	}

	if err := godotenv.Load(); err != nil {
		log.Error("No ENV Found", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(err.Error())
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Error(err.Error())
	}

	if err := db.AutoMigrate(&models.File{}); err != nil {
		log.Error(err.Error())
	}

	r := gin.Default()
	r.Use(middleware.RequestId(), middleware.AccessLog(log), gin.Recovery())
	r.MaxMultipartMemory = cfg.MaxFileSize

	repo := repositories.NewPgFileRepository(db)
	store := storage.NewLocalStorage(cfg.TempStoragePath)
	service := services.NewFileService(repo, store, log, cfg.MaxFileSize, allowedMimeTypes, allowedExtensions)
	handler := handlers.NewFileHandler(service)

	r.GET("/health", handlers.HealthCheck)
	r.POST("/files/upload", handler.Upload)
	r.DELETE("/files/delete/:id", handler.Delete)
	r.GET("/files/:id", handler.GetById)
	r.GET("/files/:id/download", handler.Download)

	uploadWorker := workers.NewUploadWorker(repo, log)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	appCtx, cancel := context.WithCancel(context.Background())
	go uploadWorker.Start(appCtx)

	<-quit
	log.Info("shutdown signal received")
	cancel()
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(shutDownCtx); err != nil {
		log.Error("server shutdown failed", "error", err)
	}
	log.Info("server exited gracefully")
}
