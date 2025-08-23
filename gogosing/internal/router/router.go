package router

import (
	bookHandler "gogosing/internal/handler/book"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateRouter(handler *bookHandler.BookHandler) *mux.Router {
	router := mux.NewRouter()

	//Add middleware
	//router.Use()

	createBookRouter(router, handler)

	return router
}

func createBookRouter(router *mux.Router, handler *bookHandler.BookHandler) {
	router.HandleFunc("/books", handler.GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/books/{id}", handler.GetBook).Methods(http.MethodGet)
	router.HandleFunc("/books", handler.CreateBook).Methods(http.MethodPost)
	router.HandleFunc("/books/{id}", handler.UpdateBook).Methods(http.MethodPut)
	router.HandleFunc("/books/{id}", handler.DeleteBook).Methods(http.MethodDelete)
}
