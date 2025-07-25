package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	bookHandler "gogosing/internal/handler/book"
	"gogosing/internal/router"
	bookStore "gogosing/internal/store/book"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	loadEnv()
	dbBookStore := bookStore.NewMySQLBookStore()

	defer func() {
		if err := dbBookStore.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	handler := bookHandler.NewBookHandler(dbBookStore)

	serverRouter := router.CreateRouter(handler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", serverRouter))
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
