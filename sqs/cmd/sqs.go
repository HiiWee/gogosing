package main

import (
	"context"
	"log"
	"sqs-example/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	ctx := context.Background()
	a := app.New(ctx)

	a.Run(ctx)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
