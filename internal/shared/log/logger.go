package log

import (
	"io"
	"log/slog"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
)

func NewLogger(cfg *config.Config, out io.Writer) *slog.Logger {
	level := slog.LevelInfo
	if cfg.AppEnviroment == "debug" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch cfg.AppEnviroment {
	case "production":
		handler = slog.NewJSONHandler(out, opts)
	default:
		handler = slog.NewTextHandler(out, opts)
	}

	return slog.New(handler)
}
