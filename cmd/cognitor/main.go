package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kernelstub/cognitor/internal/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, os.Args[1:]); err != nil {
		slog.Error("cognitor failed", "error", err)
		os.Exit(1)
	}
}
