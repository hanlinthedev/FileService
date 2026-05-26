package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	TempStoragePath      string
	PermanentStoragePath string
	MaxFileSize          int64
	AllowedExtension     []string
}

func LoadConfig() (*Config, error) {

	maxFileSize, _ := strconv.ParseInt(os.Getenv("MAX_FILE_SIZE"), 10, 64)

	config := &Config{
		Port:                 os.Getenv("PORT"),
		DBHost:               os.Getenv("DB_HOST"),
		DBPort:               os.Getenv("DB_PORT"),
		DBUser:               os.Getenv("DB_USER"),
		DBPassword:           os.Getenv("DB_PASSWORD"),
		DBName:               os.Getenv("DB_NAME"),
		TempStoragePath:      os.Getenv("TEMP_STORAGE_PATH"),
		PermanentStoragePath: os.Getenv("PERMANENT_STORAGE_PATH"),
		MaxFileSize:          maxFileSize,
	}

	return config, nil
}
