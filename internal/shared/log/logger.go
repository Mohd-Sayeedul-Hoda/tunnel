package log

import (
	"io"
	"log/slog"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
)

func NewLogger(cfg *config.Config, out io.Writer) (logger *slog.Logger) {

	minLevel := slog.LevelInfo
	if cfg.AppEnviroment == "debug" {
		minLevel = slog.LevelDebug
	}

	opts := slog.HandlerOptions{
		Level: minLevel,
	}

	var handler slog.Handler
	handler = slog.NewTextHandler(out, &opts)

	if cfg.AppEnviroment == "production" {
		handler = slog.NewJSONHandler(out, &opts)
	}

	logger = slog.New(handler)
	return
}
