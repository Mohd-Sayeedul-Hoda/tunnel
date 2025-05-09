package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
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

	cfg, err := config.InitalizeConfig(getenv, args)
	if err != nil {
		return err
	}

	log := log.NewLogger(cfg, w)
	log.Info("first info", "key", "value", "num", 1)

	return nil
}
