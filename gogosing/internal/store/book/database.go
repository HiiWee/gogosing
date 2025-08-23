package book

import (
	"database/sql"
	"errors"
	"fmt"
	"gogosing/internal/model"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLBookStore struct {
	db *sql.DB
}

func NewMySQLBookStore() *MySQLBookStore {
	store := MySQLBookStore{}
	store.connectDB()

	return &store
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

	if err := store.db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
}

func (store *MySQLBookStore) Close() error {
	return store.db.Close()
}

func (store *MySQLBookStore) GetBooks() []model.Book {
	books := make([]model.Book, 0)

	// Execute query
	rows, err := store.db.Query("SELECT id, title, author FROM books")
	if err != nil {
		log.Fatal("Error querying books:", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal("Error closing rows:", err)
		}
	}(rows)

	// Iterate over result set
	for rows.Next() {
		var book model.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			log.Fatal("Error scanning book row:", err)
		}
		books = append(books, book)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		log.Fatal("Row iteration error:", err)
	}

	return books
}

func (store *MySQLBookStore) GetBook(ID string) (model.Book, bool) {
	row := store.db.QueryRow("SELECT id, title, author FROM books WHERE id = ?", ID)

	book := model.Book{}
	err := row.Scan(&book.ID, &book.Title, &book.Author)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Book{}, false
		}
	}

	return book, true
}

func (store *MySQLBookStore) CreateBook(requestedBook model.Book) error {
	// Insert new book
	_, err := store.db.Exec(
		"INSERT INTO books (title, author) VALUES (?, ?)",
		requestedBook.Title,
		requestedBook.Author,
	)

	if err != nil {
		return err
	}
	return nil
}

func (store *MySQLBookStore) UpdateBook(ID string, requestedBook model.Book) error {
	existingBook, found := store.GetBook(ID)
	if !found {
		return errors.New("book does not exist")
	}

	title := existingBook.Title
	if requestedBook.Title != "" {
		title = requestedBook.Title
	}
	author := existingBook.Author
	if requestedBook.Author != "" {
		author = requestedBook.Author
	}

	_, err := store.db.Exec(
		"UPDATE books SET title = ?, author = ? WHERE id = ?",
		title,
		author,
		ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (store *MySQLBookStore) DeleteBook(ID string) error {
	result, err := store.db.Exec("DELETE FROM books WHERE id = ?", ID)

	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("book does not exist")
	}
	if err != nil {
		return err
	}

	return nil
}
