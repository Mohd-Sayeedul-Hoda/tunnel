package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/db"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/log"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	godotenv.Load(".env")
	var getenv func(string) string
	getenv = func(key string) string {
		return os.Getenv(key)
	}

	err := run(ctx, getenv, os.Args, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error while starting the application: %s\n", err)
		os.Exit(1)
	}

}

func run(ctx context.Context, getenv func(string) string, args []string, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := config.InitializeConfig(getenv, args)
	if err != nil {
		return err
	}

	slog.SetDefault(log.NewLogger(cfg, w))
	slog.Error("error", slog.String("path", "api/v1"))

	_, err = db.OpenDB(ctx, cfg)
	if err != nil {
		return err
	}

	slog.Info("database connected")

	return nil
}
