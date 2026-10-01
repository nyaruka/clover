package runtime

import (
	"errors"
	"log/slog"
	"strings"
)

// Config is our top level configuration object
type Config struct {
	DB       string     `help:"the connection string for our database"`
	LogLevel slog.Level `help:"the log level, one of error, warn, info, debug"`
	Version  string     `help:"the version being run"`
	Password string     `help:"the password for the admin user (required)"`
	Address  string     `help:"the address clover will listen on, empty means all interfaces"`
	Port     int        `help:"the port clover will listen on"`
}

// NewDefaultConfig returns a new default configuration object
func NewDefaultConfig() *Config {
	return &Config{
		DB:       "postgres://clover_test:temba@localhost/clover_test?sslmode=disable",
		LogLevel: slog.LevelInfo,
		Address:  "",
		Port:     8060,
		Version:  "Dev",
	}
}

// Parse validates the config
func (c *Config) Parse() error {
	// require an admin password so the admin interface is never exposed with a default
	if c.Password == "" {
		return errors.New("no admin password set")
	}

	// nothing will work right unless the connection's timezone is a constant UTC, which we force
	if strings.Contains(c.DB, "TimeZone") {
		return errors.New("db connection string can't specify a timezone, clover always uses UTC")
	}

	return nil
}
