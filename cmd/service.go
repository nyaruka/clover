package cmd

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/lib/pq"
	clover "github.com/nyaruka/clover/v26"
	"github.com/nyaruka/clover/v26/runtime"
)

// Service starts the clover service, blocks until a termination signal is received, then stops it. The config must
// already be loaded, e.g. with LoadConfig. All logging is sent to the given handler, e.g. LogHandler(), whose level
// is set from the config.
func Service(cfg *runtime.Config, version, date string, logHandler slog.Handler) error {
	// a build without a stamped version keeps the configured one
	if version != "Dev" {
		cfg.Version = version
	}

	// configure our logger
	logLevel.Set(cfg.LogLevel)
	slog.SetDefault(slog.New(logHandler))

	log := slog.With("comp", "main")
	log.Info("starting clover", "version", version, "released", date)

	// force our DB connection to be in UTC
	if strings.Contains(cfg.DB, "?") {
		cfg.DB += "&TimeZone=UTC"
	} else {
		cfg.DB += "?TimeZone=UTC"
	}

	var templateFS http.FileSystem
	if cfg.Version == "Dev" {
		templateFS = http.Dir("static")
	} else {
		templateFS = clover.Assets()
	}

	rt, err := runtime.NewRuntime(cfg)
	if err != nil {
		return err
	}

	srv := clover.NewServer(rt, templateFS)
	if err := srv.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Info("stopping clover", "signal", <-ch)

	return srv.Stop()
}
