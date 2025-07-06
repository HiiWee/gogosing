package book

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"gogosing/internal/model"
	"log"
	"os"
)

type MySQLBookStore struct {
	db *sql.DB
}

func NewMySQLBookStore() *MySQLBookStore {
	loadEnv()
	store := MySQLBookStore{}
	store.connectDB()

	return &store
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (store *MySQLBookStore) connectDB() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	store.db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.db.Close()

	// 연결 확인
	if err := store.db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
}

func (store *MySQLBookStore) GetBooks() []model.Book {
	return nil
}

func (store *MySQLBookStore) GetBook(ID string) (model.Book, bool) {
	return model.Book{}, false
}

func (store *MySQLBookStore) CreateBook(requestedBook model.Book) error {
	return nil
}

func (store *MySQLBookStore) UpdateBook(ID string, requestedBook model.Book) error {
	return nil
}

func (store *MySQLBookStore) DeleteBook(ID string) error {
	return nil
}
