package repositories

import (
	"github.com/hanlinthedev/file-service/internal/models"
	"gorm.io/gorm"
)

type pgFileRepository struct {
	db *gorm.DB
}

func NewPgFileRepository(db *gorm.DB) FileRepository {
	return &pgFileRepository{
		db: db,
	}
}

func (r *pgFileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *pgFileRepository) FindById(id string) (*models.File, error) {
	var file models.File
	if err := r.db.First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *pgFileRepository) Delete(id string) error {
	return r.db.Delete(&models.File{}, "id = ?", id).Error
}
