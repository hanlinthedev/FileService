package models

import "time"

type File struct {
	ID            string `gorm:"type:uuid;primaryKey"`
	OriginalName  string `gorm:"not null"`
	StoredName    string `gorm:"not null"`
	MimeType      string `gorm:"not null"`
	Extension     string
	Size          int64  `gorm:"not null"`
	StoragePath   string `gorm:"not null"`
	PermanentPath string
	FileStatus    FileStatus `gorm:"type:varchar(20);not null"`
	CreatedAt     time.Time
}
