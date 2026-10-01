package runtime_test

import (
	"testing"

	"github.com/nyaruka/clover/v26/runtime"
	"github.com/stretchr/testify/assert"
)

func TestConfigParse(t *testing.T) {
	cfg := runtime.NewDefaultConfig()
	assert.EqualError(t, cfg.Parse(), "no admin password set")

	cfg.Password = "sesame"
	assert.NoError(t, cfg.Parse())

	cfg.DB += "&TimeZone=America/New_York"
	assert.EqualError(t, cfg.Parse(), "db connection string can't specify a timezone, clover always uses UTC")
}
