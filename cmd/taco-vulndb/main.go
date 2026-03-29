package main

import (
	"log/slog"
	"os"

	"github.com/Joeri-Abbo/taco-vulndb/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}
