package http

import (
	"net/http"

	"bookloop.net/internal/middleware"
	"bookloop.net/internal/users"
	"github.com/julienschmidt/httprouter"
)

func RegisterUserRoutes(router *httprouter.Router, h users.Handlers, mw *middleware.MiddlewareManager) {
	router.HandlerFunc(http.MethodGet, "/v1/users", h.Register())
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", h.Activate())
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", h.Login())
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", h.CreateActivationToken())
}
