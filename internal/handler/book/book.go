package book

import (
	"gogosing/internal/model"
	bookStore "gogosing/internal/store/book"

	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type BookHandler struct {
	Store bookStore.BookStore
}

func NewBookHandler(store bookStore.BookStore) *BookHandler {
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

	var updatedBook model.Book // 요청 본문에서 업데이트할 필드만 받을 수 있도록
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
