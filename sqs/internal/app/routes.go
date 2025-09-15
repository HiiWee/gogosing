package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

type App struct {
	lgr  *slog.Logger
	r    *chi.Mux
	stop chan os.Signal
}

func (a *App) New() *App {
	app := &App{
		lgr: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})),
		r:    chi.NewRouter(),
		stop: make(chan os.Signal, 1),
	}
	slog.SetDefault(app.lgr)
	signal.Notify(app.stop, syscall.SIGINT, syscall.SIGTERM)

	return app
}
