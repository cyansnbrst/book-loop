package http

import (
	"net/http"

	"bookloop.net/internal/books"
	"bookloop.net/internal/middleware"
	"github.com/julienschmidt/httprouter"
)

func RegisterBookRoutes(router *httprouter.Router, h books.Handlers, mw *middleware.MiddlewareManager) {
	router.HandlerFunc(http.MethodGet, "/v1/books", mw.RequirePermission("books:read", h.List()))
	router.HandlerFunc(http.MethodPost, "/v1/books", mw.RequirePermission("books:write", h.Create()))
	router.HandlerFunc(http.MethodGet, "/v1/books/:id", mw.RequirePermission("books:read", h.Get()))
	router.HandlerFunc(http.MethodPatch, "/v1/books/:id", mw.RequirePermission("books:write", h.Update()))
	router.HandlerFunc(http.MethodDelete, "/v1/books/:id", mw.RequirePermission("books:write", h.Delete()))
}
