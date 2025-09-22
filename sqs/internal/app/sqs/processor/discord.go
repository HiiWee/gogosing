package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sqs-example/internal/app/util"
)

type Processor struct {
	webhookURL string
}

type SendingMessage interface {
	GetMessage() string
	GetFrom() string
}

func NewProcessor() *Processor {
	return &Processor{
		webhookURL: util.MustEnv("DISCORD_WEBHOOK_URL"),
	}
}

func (p *Processor) ProcessMessage(ctx context.Context, m SendingMessage) error {
	fmt.Println("message is processed by " + m.GetFrom())
	discordPayload := map[string]string{
		"content": m.GetMessage(),
		"from":    m.GetFrom(),
	}

	b, _ := json.Marshal(discordPayload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, p.webhookURL, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("error processing event: %d %s", resp.StatusCode, string(body))
	}

	log.Printf("posted to Discord: %s", m.GetMessage())

	return nil
}
