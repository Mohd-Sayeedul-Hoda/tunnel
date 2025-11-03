package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"

	natserver "github.com/Mohd-Sayeedul-Hoda/tunnel/internal/nat-server"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/db"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/log"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	var getenv func(string) string
	getenv = func(s string) string {
		return os.Getenv(s)
	}
	err := run(context.Background(), getenv, os.Args, os.Stdout)
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

	pgPool, err := db.OpenPostgresConn(ctx, cfg)

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("tcp server running")
		err := natserver.ListenAndServer(ctx, w, cfg)
		serverErrors <- err
	}()

	select {
	case <-ctx.Done():
		slog.Info("nat server shutdown initiated", slog.String("reason", "context cancelled"))
	case err := <-serverErrors:
		return err
	}

	return nil
}
