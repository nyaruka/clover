package cmd

import (
	"log/slog"
	"testing"

	"github.com/nyaruka/clover/v26/runtime"
	"github.com/nyaruka/ezconf"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg := runtime.NewDefaultConfig()
	err := loadConfig(cfg, []string{"-password", "sesame", "-log-level", "debug", "-port", "9000"})
	assert.NoError(t, err)
	assert.Equal(t, "sesame", cfg.Password)
	assert.Equal(t, slog.LevelDebug, cfg.LogLevel)
	assert.Equal(t, 9000, cfg.Port)

	// config which loads but doesn't validate
	err = loadConfig(runtime.NewDefaultConfig(), []string{})
	assert.EqualError(t, err, "invalid config: no admin password set")

	err = loadConfig(runtime.NewDefaultConfig(), []string{"-help"})
	assert.ErrorIs(t, err, ezconf.ErrHelp)
}
