package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("DB_HOST", "localhost")
	defer os.Clearenv()

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "localhost", cfg.DBHost)
}
