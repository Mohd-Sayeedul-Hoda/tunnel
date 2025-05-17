package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api"
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

	_, err = db.OpenDB(ctx, cfg)
	if err != nil {
		return err
	}

	slog.Info("database connection pool establish")

	handler := api.NewHTTPServer(cfg)
	slog.Debug("handler created")

	httpServer := http.Server{
		Addr:    net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		Handler: handler,
	}

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("http server running",
			slog.String("host", cfg.Server.Host),
			slog.Int("port", cfg.Server.Port),
			slog.String("app-env", cfg.AppEnviroment),
			slog.Int("app-version", cfg.AppVersion),
		)

		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			slog.Error("error while starting http server", slog.Any("err", err))
			serverErrors <- err
		}
	}()

	var wg sync.WaitGroup

	select {
	case <-ctx.Done():
		slog.Info("shutdown initiated", slog.String("reason", "context cancelled"))
	case err := <-serverErrors:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = httpServer.Shutdown(shutdownCtx)
	if err != nil {
		return err
	}

	slog.Info("closing http server", slog.String("addrs", httpServer.Addr))

	wg.Wait()

	slog.Info("http server stop", slog.String("addrs", httpServer.Addr))

	return nil
}
