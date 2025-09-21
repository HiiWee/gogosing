package producer

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sqs-example/internal/app/util"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	lgr      *slog.Logger
	r        *chi.Mux
	stop     chan os.Signal
	queueURL string
	producer *Producer
}

func New(ctx context.Context) *App {
	queueURL := util.MustEnv("SQS_QUEUE_URL")
	region := util.MustEnv("AWS_REGION")

	app := &App{
		lgr: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})),
		r:        chi.NewRouter(),
		stop:     make(chan os.Signal, 1),
		queueURL: queueURL,
	}
	slog.SetDefault(app.lgr)
	signal.Notify(app.stop, syscall.SIGINT, syscall.SIGTERM)

	client := NewClient(ctx, region)
	app.producer = NewProducer(client)
	app.registerRoutes(ctx)

	return app
}

func (a *App) Run(ctx context.Context) {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.TimeoutHandler(a.r, 1*time.Second, "timed out"),
	}

	a.lgr.Info("starting server")
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.lgr.Error("listen err", "error", err)
		}
	}()

	<-a.stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		a.lgr.Error("app setup failed: graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	a.lgr.Info("server shutdown complete")
}
