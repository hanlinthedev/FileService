package storage

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	rootPath string
}

func NewLocalStorage(rootPath string) Storage {
	return &LocalStorage{rootPath: rootPath}
}

func (s *LocalStorage) Save(file multipart.File, path string) error {
	fullPath := filepath.Join(s.rootPath, path)

	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer func() {
		if err := dst.Close(); err != nil {
			log.Println(err)
		}
	}()
	_, err = io.Copy(dst, file)
	return err
}

func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(s.rootPath, path)
	return os.Remove(fullPath)
}

func (s *LocalStorage) GetFullPath(path string) string {
	return filepath.Join(s.rootPath, path)
}

func (s *LocalStorage) Exists(path string) bool {
	full := s.GetFullPath(path)
	_, err := os.Stat(full)
	return err == nil
}
