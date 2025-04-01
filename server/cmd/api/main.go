package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	godotenv.Load(".env")
	var getenv func(string) string
	getenv = func(key string) string {
		return os.Getenv(key)
	}

	err := run(ctx, getenv, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error while starting the application: %s\n", err)
		os.Exit(1)
	}

}

func run(ctx context.Context, getenv func(string) string, w io.Writer) error {
	return fmt.Errorf("todo")
}
