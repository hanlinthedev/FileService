package errors

import "errors"

var (
	ErrEmptyFile       = errors.New("empty file")
	ErrFileTooLarge    = errors.New("file too large")
	ErrUnsupportedMime = errors.New("unsupported mime type")
	ErrUnsupportedExt  = errors.New("unsupported file extension")
	ErrFileNotFound    = errors.New("file not found")
)
