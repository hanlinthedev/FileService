package services

import (
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
}

type fileService struct {
	repo    repositories.FileRepository
	storage storage.Storage
}

func NewFileService(repo repositories.FileRepository, storage storage.Storage) FileService {
	return &fileService{
		repo:    repo,
		storage: storage,
	}
}

func (s *fileService) Upload(file multipart.File, header *multipart.FileHeader) (*models.File, error) {
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
	panic("implement me")
}

func (s *fileService) Delete(path string) error {
	panic("implement me")
}
