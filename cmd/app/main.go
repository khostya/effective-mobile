package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/khostya/effective-mobile/internal/app"
	"github.com/khostya/effective-mobile/internal/config"
	"github.com/khostya/effective-mobile/pkg/log/sl"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("error", sl.Err(err))
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     parseLogLevel(cfg.Env),
		AddSource: true,
	})))

	if err := app.Run(ctx, cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("error", sl.Err(err))
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

const (
	infoLog  = "info"
	debugLog = "debug"
)

func parseLogLevel(level string) slog.Level {
	switch level {
	case infoLog:
		return slog.LevelInfo
	case debugLog:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
