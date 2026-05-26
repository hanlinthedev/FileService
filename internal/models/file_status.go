package models

type FileStatus string

const (
	FileStatusStaged FileStatus = "STAGED"
	FileStatusReady  FileStatus = "READY"
	FileStatusFailed FileStatus = "FAILED"
)
