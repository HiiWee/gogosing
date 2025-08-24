package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("load AWS config: %v", err)
	}
	client := sqs.NewFromConfig(awsCfg)

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		if msg == "" {
			http.Error(w, "missing msg query param", http.StatusBadRequest)
			return
		}
		from := r.URL.Query().Get("from")

		payload := Msg{
			Text:   msg,
			SentAt: time.Now().UTC(),
			From:   from,
		}
		body, _ := json.Marshal(payload)

		_, err := client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    &queueURL,
			MessageBody: awsString(string(body)),
			// Optional attributes (e.g., for filtering)
			// MessageAttributes: map[string]types.MessageAttributeValue{ ... }
		})
		if err != nil {
			log.Printf("SendMessage error: %v", err)
			http.Error(w, "failed to enqueue", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("queued\n"))
	})

	log.Println("Producer listening on :8080  ->  GET /send?msg=hello&from=hoseok")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func awsString(s string) *string { return &s }
