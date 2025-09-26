package main

import (
	"sqs-example/cmd"
	"sqs-example/internal/app/sqs/listener"
)

func main() {
	cmd.LoadEnv()
	a := listener.New()
	a.Run()
}
