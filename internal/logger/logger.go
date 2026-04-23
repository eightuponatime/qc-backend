package logger

import (
	"log/slog"
	"os"
	"qc/config"
)

func Setup(cfg *config.Config) *slog.Logger {
	options := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: cfg.Env == "development",
	}

	var handler slog.Handler
	if cfg.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	logger := slog.New(handler).With(
		slog.String("service", "qc-backend"),
		slog.String("env", cfg.Env),
	)

	slog.SetDefault(logger)
	return logger
}
