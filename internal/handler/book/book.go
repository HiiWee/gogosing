package book

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gogosing/internal/model"
	"net/http"
)

type BookStore interface {
	GetBooks() []model.Book
	GetBook(ID string) (model.Book, bool)
	CreateBook(requestedBook model.Book) error
	UpdateBook(ID string, requestedBook model.Book) error
	DeleteBook(ID string) error
}

type BookHandler struct {
	Store BookStore
}

func NewBookHandler(store BookStore) *BookHandler {
	return &BookHandler{Store: store}
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books := h.Store.GetBooks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	book, ok := h.Store.GetBook(params["id"])

	if ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
		return
	}
	http.NotFound(w, r)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	err := h.Store.CreateBook(book)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var updatedBook model.Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)

	err = h.Store.UpdateBook(params["id"], updatedBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	err := h.Store.DeleteBook(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
