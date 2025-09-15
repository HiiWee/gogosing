package sqs

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client struct {
	client *sqs.Client
}

func New(ctx context.Context) *Client {
	region := mustEnv("AWS_REGION")

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

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}
