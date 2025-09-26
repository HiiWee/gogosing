package producer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type PublishEvent struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

type Producer struct {
	sender Sender
}

type Sender interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

func NewProducer(sender Sender) *Producer {
	return &Producer{
		sender: sender,
	}
}

func (p *Producer) SendMessage(ctx context.Context, event *PublishEvent, queueURL string) error {
	payload, err := json.Marshal(event)

	if err != nil {
		return err
	}

	_, err = p.sender.SendMessage(ctx, &sqs.SendMessageInput{MessageBody: aws.String(string(payload)), QueueUrl: aws.String(queueURL)})

	if err != nil {
		slog.Error("failed to send message", "error", err)
		return err
	}
	return nil
}
