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
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/log"
	"github.com/hashicorp/yamux"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	var getenv func(string) string
	getenv = func(s string) string {
		return os.Getenv(s)
	}
	err := run(context.Background(), getenv, os.Args, os.Stdin)
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

	listner, err := net.Listen("tcp", cfg.NatServer.Host+":"+strconv.Itoa(cfg.NatServer.Port))
	if err != nil {
		return err
	}
	slog.Info("tcp server started")

	Q serverErrors := make(chan error, 1)
		conn, err := listner.Accept()
		if err != nil {
			serverErrors <- err
		}

		natserver.HandleConnection(cfg, conn)
	}()

	select {
	case <-ctx.Done():
		slog.Info("nat server shutdown initiated", slog.String("reason", "context cancelled"))
		return listner.Close()
	case err := <-serverErrors:
		return err
	}

}
