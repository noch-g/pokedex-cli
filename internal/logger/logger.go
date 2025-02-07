package logger

import (
	"flag"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	opts := &slog.HandlerOptions{}
	if *debug {
		opts.Level = slog.LevelDebug
	} else {
		opts.Level = slog.LevelInfo
	}

	Logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(Logger)
}

// Fonctions raccourcies
func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}
