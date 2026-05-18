package database

import (
	"testing"

	"github.com/hanlinthedev/file-service/internal/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConnection_Integration(t *testing.T) {
	_ = godotenv.Load("../../.env.test")
	cfg, _ := config.LoadConfig()

	db, err := NewConnection(cfg)

	require.NoError(t, err)
	require.NotNil(t, db)

	type TestModel struct {
		ID uint64 `gorm:"primary_key"`
	}
	err = db.AutoMigrate(&TestModel{})

	assert.NoError(t, err)
}
