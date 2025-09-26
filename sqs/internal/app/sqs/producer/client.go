package producer

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client struct {
	client *sqs.Client
}

func NewClient(region string) *Client {
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(awsCfg)

	return &Client{
		client: client,
	}
}

func (c *Client) SendMessage(ctx context.Context, params *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return c.client.SendMessage(ctx, params)
}
