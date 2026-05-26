package logger

import (
	"log/slog"
	"os"
	"time"
)

type NewLog struct {
	requestId string
	userId    string
	status    string
	logAt     time.Time
}

func AppLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}
