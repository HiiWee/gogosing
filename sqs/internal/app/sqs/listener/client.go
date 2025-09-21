package listener

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client struct {
	client *sqs.Client
}

func NewClient(ctx context.Context, region string) *Client {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(awsCfg)

	return &Client{
		client: client,
	}
}

func (c *Client) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return c.client.ReceiveMessage(ctx, params)
}

func (c *Client) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return c.client.DeleteMessage(ctx, params)
}
