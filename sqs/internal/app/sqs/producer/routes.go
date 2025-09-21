package producer

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
)

func (a *App) registerRoutes(ctx context.Context) {
	a.r.Post("/send", a.buildMessageSending(ctx))
}

func (a *App) buildMessageSending(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.lgr.Info("buildMessageSending")
		var event PublishEvent
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
