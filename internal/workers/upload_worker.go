package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/hanlinthedev/file-service/internal/models"
	"github.com/hanlinthedev/file-service/internal/repositories"
)

type UploadWorker struct {
	repo repositories.FileRepository
	log  *slog.Logger
}

func NewUploadWorker(
	repo repositories.FileRepository, log *slog.Logger) *UploadWorker {
	return &UploadWorker{
		repo: repo,
		log:  log,
	}
}

func (w *UploadWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	w.log.Info("upload worker started")
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Info("upload worker stopped")
			return
		case <-ticker.C:
			files, err := w.repo.FindByStatus(ctx, models.FileStatusStaged)
			if err != nil {
				w.log.Error("failed to fetch staged data", "error", err.Error())
				continue
			}
			for _, file := range files {
				err := w.repo.UpdateStatus(ctx, file.ID, models.FileStatusProcessing)
				if err != nil {
					w.log.Error("failed to mark file processing", "file_id", file.ID, "error", err)
					continue
				}
				w.log.Info("processing file", "file_id ", file.ID)
				time.Sleep(5 * time.Second)
				err = w.repo.UpdateStatus(ctx, file.ID, models.FileStatusReady)
				if err != nil {
					w.log.Error("failed to mark file ready", "file_id", file.ID, "error", err)
					continue
				}
				w.log.Info("file processed successfully", "file_id", file.ID)
			}
			w.log.Info("staged files found", "count", len(files))
		}
	}
}
