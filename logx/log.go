package logx

import (
	"log/slog"
	"os"
)

var Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
