package repositories

import (
	"context"

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

func (r *pgFileRepository) Create(ctx context.Context, file *models.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *pgFileRepository) FindById(id string) (*models.File, error) {
	var file models.File
	if err := r.db.First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *pgFileRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.File{}, "id = ?", id).Error
}

func (r *pgFileRepository) FindByStatus(ctx context.Context, status models.FileStatus) ([]models.File, error) {
	var files []models.File
	if err := r.db.WithContext(ctx).Where("file_status = ?", status).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *pgFileRepository) UpdateStatus(ctx context.Context, id string, status models.FileStatus) error {
	return r.db.Model(&models.File{}).WithContext(ctx).Where("id = ?", id).Update("file_status", status).Error
}
