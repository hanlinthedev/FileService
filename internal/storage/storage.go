package storage

import "mime/multipart"

type Storage interface {
	Save(file multipart.File, path string) error
	Delete(path string) error
}
