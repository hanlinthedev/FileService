package services

import (
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	customeErrors "github.com/hanlinthedev/file-service/internal/errors"
	"github.com/hanlinthedev/file-service/internal/models"
	"github.com/hanlinthedev/file-service/internal/repositories"
	"github.com/hanlinthedev/file-service/internal/storage"
)

type FileService interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*models.File, error)
	GetById(id string) (*models.File, error)
	Delete(ctx context.Context, id string) error
	GetStoragePath(id string) (string, string, error)
}

type fileService struct {
	repo    repositories.FileRepository
	storage storage.Storage
	log     *slog.Logger

	maxFileSize      int64
	allowedMimeType  map[string]bool
	allowedExtension map[string]bool
}

func NewFileService(repo repositories.FileRepository, storage storage.Storage, log *slog.Logger, maxFileSize int64, allowedMimeType map[string]bool, allowedExtension map[string]bool) FileService {
	return &fileService{
		repo:             repo,
		storage:          storage,
		log:              log,
		maxFileSize:      maxFileSize,
		allowedMimeType:  allowedMimeType,
		allowedExtension: allowedExtension,
	}
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

func (s *fileService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*models.File, error) {

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

	if err := s.validateFile(header, mimeType, ext); err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	if err := s.storage.Save(file, relativePath); err != nil {
		s.log.Error("error saving file", err)
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
		FileStatus:   models.FileStatusStaged,
	}

	err = s.repo.Create(ctx, model)
	if err != nil {
		_ = s.storage.Delete(relativePath)
		s.log.Error("error saving record to database", err)
		return nil, err
	}

	return model, nil
}

func (s *fileService) GetById(id string) (*models.File, error) {
	return s.repo.FindById(id)
}

func (s *fileService) Delete(ctx context.Context, id string) error {
	file, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, file.ID); err != nil {
		return err
	}

	if err := s.storage.Delete(file.StoragePath); err != nil {
		_ = s.repo.Create(ctx, file).Error()
		return err
	}

	return nil
}

func (s *fileService) validateFile(header *multipart.FileHeader, mimeType string, ext string) error {
	if header.Size <= 0 {
		return customeErrors.ErrEmptyFile
	}
	if header.Size > s.maxFileSize {
		return customeErrors.ErrFileTooLarge
	}
	if !s.allowedMimeType[mimeType] {
		return customeErrors.ErrUnsupportedMime
	}

	if !s.allowedExtension[ext] {
		return customeErrors.ErrUnsupportedExt
	}
	return nil
}
