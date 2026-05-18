package repositories

import "github.com/hanlinthedev/file-service/internal/models"

type FileRepository interface {
	Create(file *models.File) error
	FindById(id string) (*models.File, error)
	Delete(id string) error
}
