package listener

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sqs-example/internal/app/sqs/processor"
	"sqs-example/internal/app/util"
	"syscall"
)

type App struct {
	lgr      *slog.Logger
	stop     chan os.Signal
	queueURL string
	listener *Listener
}

func New(ctx context.Context) *App {
	queueURL := util.MustEnv("SQS_QUEUE_URL")
	region := util.MustEnv("AWS_REGION")

	app := &App{
		lgr: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})),
		stop:     make(chan os.Signal, 1),
		queueURL: queueURL,
	}
	slog.SetDefault(app.lgr)
	signal.Notify(app.stop, syscall.SIGINT, syscall.SIGTERM)

	c := NewClient(ctx, region)
	p := processor.NewProcessor()
	app.listener = NewListener(c, p)

	return app
}

func (a *App) Run(ctx context.Context) {
	withCancel, cancel := context.WithCancel(ctx)

	a.listener.Listen(withCancel, a.queueURL)

	a.lgr.Info("starting listener")

	<-a.stop

	defer cancel()
	a.lgr.Info("listener shutdown complete")
}
