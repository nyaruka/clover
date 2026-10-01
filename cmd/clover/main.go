package main

import (
	"github.com/nyaruka/clover/v26/cmd"
	"github.com/nyaruka/clover/v26/runtime"
)

var (
	// overridden at build time via -ldflags "-X main.version=... -X main.date=..."
	version = "Dev"
	date    = "unknown"
)

func main() {
	cfg := runtime.NewDefaultConfig()
	cmd.LoadConfig(cfg)
	cmd.Run(cmd.Service(cfg, version, date, cmd.LogHandler()))
}
