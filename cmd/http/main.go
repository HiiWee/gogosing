package main

import (
	"fmt"
	bookHandler "gogosing/internal/handler/book"
	"gogosing/internal/router"
	bookStore "gogosing/internal/store/book"
	"os"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	loadEnv()
	connectDB()

	dbBookStore := bookStore.NewMySQLBookStore()
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

func connectDB() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 연결 확인
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
}
