package main

import (
	"context"
	"errors"
	"github.com/khostya/effective-mobile/internal/app"
	"github.com/khostya/effective-mobile/internal/config"
	"log"
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
		log.Fatalln(err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     parseLogLevel(cfg.Env),
		AddSource: true,
	})))

	if err := app.Run(ctx, cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
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
