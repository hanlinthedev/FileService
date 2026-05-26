package repositories

import (
	"context"

	"github.com/hanlinthedev/file-service/internal/models"
)

type FileRepository interface {
	Create(ctx context.Context, file *models.File) error
	FindById(id string) (*models.File, error)
	Delete(ctx context.Context, id string) error
}
