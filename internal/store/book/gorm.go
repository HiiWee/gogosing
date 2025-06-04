package book

import (
	"gogosing/internal/model"
	"gorm.io/gorm"
)

type GormBookStore struct {
	db *gorm.DB
}

func NewGormBookStore(db *gorm.DB) *GormBookStore {
	// auto migrate table
	db.AutoMigrate(&model.Book{})
	return &GormBookStore{db: db}
}

func (store *GormBookStore) GetBooks() []model.Book {
	var books []model.Book
	store.db.Find(&books)
	return books
}

func (store *GormBookStore) GetBook(ID string) (model.Book, bool) {
	var book model.Book
	result := store.db.First(&book, "id = ?", ID)
	if result.Error != nil {
		return model.Book{}, false
	}
	return book, true
}

func (store *GormBookStore) CreateBook(requestedBook model.Book) error {
	return store.db.Create(&requestedBook).Error
}

func (store *GormBookStore) UpdateBook(ID string, requestedBook model.Book) error {
	var book model.Book
	if err := store.db.First(&book, "id = ?", ID).Error; err != nil {
		return err
	}
	if requestedBook.Title != "" {
		book.Title = requestedBook.Title
	}
	if requestedBook.Author != "" {
		book.Author = requestedBook.Author
	}
	return store.db.Save(&book).Error
}

func (store *GormBookStore) DeleteBook(ID string) error {
	return store.db.Delete(&model.Book{}, "id = ?", ID).Error
}
