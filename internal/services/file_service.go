package services

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hanlinthedev/file-service/internal/models"
	"github.com/hanlinthedev/file-service/internal/repositories"
	"github.com/hanlinthedev/file-service/internal/storage"
)

type FileService interface {
	Upload(file multipart.File, header *multipart.FileHeader) (*models.File, error)
	GetById(id string) (*models.File, error)
	Delete(id string) error
	GetStoragePath(id string) (string, string, error)
}

type fileService struct {
	repo        repositories.FileRepository
	storage     storage.Storage
	maxFileSize int64
}

func (s *fileService) GetStoragePath(id string) (string, string, error) {
	file, err := s.repo.FindById(id)
	if err != nil {
		return "", "", err
	}

	if !s.storage.Exists(file.StoragePath) {
		return "", "", err
	}

	fullPath := s.storage.GetFullPath(file.StoragePath)

	return fullPath, file.MimeType, nil
}

func NewFileService(repo repositories.FileRepository, storage storage.Storage, maxFileSize int64) FileService {
	return &fileService{
		repo:        repo,
		storage:     storage,
		maxFileSize: maxFileSize,
	}
}

func (s *fileService) Upload(file multipart.File, header *multipart.FileHeader) (*models.File, error) {

	if header.Size <= 0 {
		return nil, errors.New("empty file")
	}
	if header.Size > s.maxFileSize {
		return nil, errors.New("file too large")
	}

	id := uuid.New().String()
	ext := filepath.Ext(header.Filename)

	storedName := id + ext
	now := time.Now()

	relativePath := filepath.Join(now.Format("2006"), now.Format("01"), storedName)

	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	mimeType := http.DetectContentType(buffer)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	if err := s.storage.Save(file, relativePath); err != nil {
		return nil, err
	}

	model := &models.File{
		ID:           id,
		OriginalName: header.Filename,
		StoredName:   storedName,
		MimeType:     mimeType,
		Extension:    ext,
		Size:         header.Size,
		StoragePath:  relativePath,
	}

	err = s.repo.Create(model)
	if err != nil {
		_ = s.storage.Delete(relativePath)
		return nil, err
	}

	return model, nil
}

func (s *fileService) GetById(id string) (*models.File, error) {
	return s.repo.FindById(id)
}

func (s *fileService) Delete(id string) error {
	file, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(file.ID); err != nil {
		return err
	}

	if err := s.storage.Delete(file.StoragePath); err != nil {
		_ = s.repo.Create(file).Error()
		return err
	}

	return nil
}
