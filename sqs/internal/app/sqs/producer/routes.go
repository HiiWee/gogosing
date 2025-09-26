package producer

import (
	"encoding/json"
	"net/http"
)

func (a *App) registerRoutes() {
	a.r.Post("/send", a.buildMessageSending())
}

func (a *App) buildMessageSending() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		a.lgr.Info("buildMessageSending")
		var event PublishEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			a.lgr.Error("error decoding json ", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := a.producer.SendMessage(ctx, &event, a.queueURL)
		if err != nil {
			a.lgr.Error("failed to send message", "error", err)
		}
	}
}
