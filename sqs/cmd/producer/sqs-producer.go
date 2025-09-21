package main

import (
	"context"
	"sqs-example/cmd"
	"sqs-example/internal/app/sqs/producer"
)

func main() {
	cmd.LoadEnv()
	ctx := context.Background()
	a := producer.New(ctx)

	a.Run(ctx)
}
