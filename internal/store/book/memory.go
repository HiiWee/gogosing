package book

import (
	"fmt"
	"gogosing/internal/model"
)

type BookStore interface {
	GetBooks() []model.Book
	GetBook(ID string) (model.Book, bool)
	CreateBook(requestedBook model.Book) error
	UpdateBook(ID string, requestedBook model.Book) error
	DeleteBook(ID string) error
}

type InMemoryBookStore struct {
	books map[string]model.Book
}

func NewInMemoryBookStore() *InMemoryBookStore {
	return &InMemoryBookStore{
		books: make(map[string]model.Book),
	}
}

func (store *InMemoryBookStore) GetBooks() []model.Book {
	bookList := make([]model.Book, 0, len(store.books)) // 전역 books 맵의 크기로 슬라이스 초기화
	for _, book := range store.books {                  // 전역 books 맵 순회
		bookList = append(bookList, book)
	}

	return bookList
}

func (store *InMemoryBookStore) GetBook(ID string) (model.Book, bool) {
	book, ok := store.books[ID]

	if ok {
		return book, true
	}

	return model.Book{}, false
}

func (store *InMemoryBookStore) CreateBook(requestedBook model.Book) error {
	if _, exists := store.books[requestedBook.ID]; exists {
		return fmt.Errorf("book with ID %s already exists", requestedBook.ID)
	}
	store.books[requestedBook.ID] = requestedBook

	return nil
}

func (store *InMemoryBookStore) UpdateBook(ID string, requestedBook model.Book) error {
	if _, exists := store.books[ID]; !exists {
		return fmt.Errorf("book with ID %s not exists", ID)
	}

	targetBook := store.books[ID]

	if requestedBook.Title != "" {
		targetBook.Title = requestedBook.Title
	}
	if requestedBook.Author != "" {
		targetBook.Author = requestedBook.Author
	}
	store.books[ID] = targetBook

	return nil
}

func (store *InMemoryBookStore) DeleteBook(ID string) error {
	if _, exists := store.books[ID]; !exists {
		return fmt.Errorf("book with ID %s does not exist", ID)
	}
	delete(store.books, ID)

	return nil
}
