package main

import (
	"context"
	"sqs-example/cmd"
	"sqs-example/internal/app/sqs/listener"
)

func main() {
	cmd.LoadEnv()
	ctx := context.Background()
	a := listener.New(ctx)

	a.Run(ctx)
}
