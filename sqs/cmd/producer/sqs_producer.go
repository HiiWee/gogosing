package main

import (
	"sqs-example/cmd"
	"sqs-example/internal/app/sqs/producer"
)

func main() {
	cmd.LoadEnv()
	a := producer.New()
	a.Run()
}
