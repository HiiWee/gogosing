package sqs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type ConsumedEvent struct {
	from    string
	message string
}
type Listener struct {
	eventChannel chan ConsumedEvent
	receiver     Receiver
}

type Receiver interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
}

func NewListener(r Receiver) *Listener {
	return &Listener{
		eventChannel: make(chan ConsumedEvent),
		receiver:     r,
	}
}

func (l *Listener) Listen(ctx context.Context, url string) {
	go l.listen(ctx, url)
	go l.process(ctx)
}

func (l *Listener) listen(ctx context.Context, url string) {
	for {
		out, err := l.receiver.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &url,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10, // long polling
			VisibilityTimeout:   30, // seconds to process before it reappears
		})
		if err != nil {
			log.Printf("ReceiveMessage error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if len(out.Messages) == 0 {
			continue
		}

		for _, msg := range out.Messages {
			var event ConsumedEvent
			err := json.Unmarshal([]byte(*msg.Body), &event)

			if err != nil {
				log.Printf("Unmarshal error: %v", err)
			}
			l.eventChannel <- event
		}
	}
}

func (l *Listener) process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-l.eventChannel:
			l.processEvent(ctx, &e)
		}
	}
}

func (l *Listener) processEvent(ctx context.Context, e *ConsumedEvent) {
	discordWebhook := mustEnv("DISCORD_WEBHOOK_URL")

	fmt.Println("message is processed by " + e.from)
	discordPayload := map[string]string{
		"content": e.message,
		"from":    e.from,
	}

	b, _ := json.Marshal(discordPayload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, discordWebhook, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error processing event: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Error processing event: %d %s", resp.StatusCode, string(body))
	}

	log.Printf("posted to Discord: %s", e.message)
}
