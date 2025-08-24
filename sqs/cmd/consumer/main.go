package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/joho/godotenv"
)

type Msg struct {
	Text   string    `json:"text"`
	SentAt time.Time `json:"sent_at"`
	From   string    `json:"from,omitempty"`
}

func main() {
	loadEnv()
	queueURL := mustEnv("SQS_QUEUE_URL")
	region := mustEnv("AWS_REGION")
	discordWebhook := mustEnv("DISCORD_WEBHOOK_URL")

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("load AWS config: %v", err)
	}
	client := sqs.NewFromConfig(awsCfg)

	log.Println("Consumer started: polling SQS and posting to Discord")

	for {
		out, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
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

		for _, m := range out.Messages {
			if err := handleMessage(ctx, client, queueURL, discordWebhook, m); err != nil {
				// Let message become visible again for retry; optionally change visibility or send to DLQ
				log.Printf("handleMessage failed (will retry later): %v", err)
			} else {
				_, _ = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &queueURL,
					ReceiptHandle: m.ReceiptHandle,
				})
			}
		}
	}
}

func handleMessage(ctx context.Context, client *sqs.Client, queueURL, discordWebhook string, m types.Message) error {
	var msg Msg
	if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
		// If bad payload, you might want to delete it to avoid poison message loops
		log.Printf("bad payload, deleting: %v", err)
		_, _ = client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &queueURL,
			ReceiptHandle: m.ReceiptHandle,
		})
		return nil
	}

	content := msg.Text
	if msg.From != "" {
		content = "**" + msg.From + "**: " + content
	}
	discordPayload := map[string]string{"content": content}

	b, _ := json.Marshal(discordPayload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, discordWebhook, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &httpError{status: resp.StatusCode, body: string(body)}
	}

	log.Printf("posted to Discord: %s", content)
	return nil
}

type httpError struct {
	status int
	body   string
}

func (e *httpError) Error() string {
	return "discord http " + http.StatusText(e.status) + ": " + e.body
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}
