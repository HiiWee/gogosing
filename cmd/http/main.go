package main

import (
	bookHandler "gogosing/internal/handler/book"
	"gogosing/internal/router"
	bookStore "gogosing/internal/store/book"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {

	connectDB()

	store := bookStore.NewInMemoryBookStore()
	handler := bookHandler.NewBookHandler(store)

	serverRouter := router.CreateRouter(handler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", serverRouter))
}

func connectDB() {
	var err error
	dsn := "hoseok:1234@tcp(localhost:3306)/gogosing"
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
