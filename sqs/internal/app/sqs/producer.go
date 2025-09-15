package sqs

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type PublishEvent struct {
	from    string
	message string
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

func (p *Producer) SendMessage(ctx context.Context, event *PublishEvent) error {
	payload, err := json.Marshal(event)

	if err != nil {
		return err
	}

	_, err = p.sender.SendMessage(ctx, &sqs.SendMessageInput{MessageBody: aws.String(string(payload))})

	if err != nil {
		log.Printf("failed to send message: %v", err)
		return err
	}
	return nil
}
