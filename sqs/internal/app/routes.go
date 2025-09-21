package app

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sqs-example/internal/app/sqs"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	lgr      *slog.Logger
	r        *chi.Mux
	stop     chan os.Signal
	queueURL string
	producer *sqs.Producer
	listener *sqs.Listener
}

func New(ctx context.Context) *App {
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		slog.Error("SQS_QUEUE_URL environment variable not set")
		return nil
	}

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

	client := sqs.NewClient(ctx)
	app.producer = sqs.NewProducer(client)
	app.listener = sqs.NewListener(client)
	app.registerRoutes(ctx)

	return app
}

func (a *App) registerRoutes(ctx context.Context) {
	a.r.Post("/send", a.buildMessageSending(ctx))
}

func (a *App) buildMessageSending(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.lgr.Info("buildMessageSending")
		var event sqs.PublishEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			a.lgr.Error("error decoding json ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := os.Getenv("SQS_QUEUE_URL")
		if url == "" {
			a.lgr.Error("error getting SQS_QUEUE_URL")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := a.producer.SendMessage(ctx, &event, a.queueURL)
		if err != nil {
			a.lgr.Error("failed to send message", err)
		}
	}
}

func (a *App) Run(ctx context.Context) {
	withCancel, cancel := context.WithCancel(ctx)

	a.listener.Listen(withCancel, a.queueURL)

	a.lgr.Info("starting listener")

	defer cancel()

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

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		a.lgr.Error("app setup failed: graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	a.lgr.Info("server shutdown complete")
}
