package models

type FileStatus string

const (
	FileStatusStaged     FileStatus = "STAGED"
	FileStatusProcessing FileStatus = "PROCESSING"
	FileStatusReady      FileStatus = "READY"
	FileStatusFailed     FileStatus = "FAILED"
)
