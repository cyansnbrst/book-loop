package http

import (
	"net/http"

	"bookloop.net/internal/books"
	"github.com/julienschmidt/httprouter"
)

func RegisterBookRoutes(router *httprouter.Router, h books.Handlers) {
	router.HandlerFunc(http.MethodGet, "/v1/books", h.List())
	router.HandlerFunc(http.MethodPost, "/v1/books", h.Create())
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", h.Get())
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", h.Update())
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", h.Delete())
}
