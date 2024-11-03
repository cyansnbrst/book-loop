package server

import (
	"net/http"

	booksHttp "bookloop.net/internal/books/delivery/http"
	booksRepository "bookloop.net/internal/books/repository"
	booksUseCase "bookloop.net/internal/books/usecase"
	"bookloop.net/pkg/error_responses"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error_responses.NotFoundResponse(w, r, s.logger)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error_responses.MethodNotAllowedResponse(w, r, s.logger)
	})

	// Init repositories
	booksRepo := booksRepository.NewBooksRepository(s.db)

	// Init useCases
	booksUC := booksUseCase.NewBooksUseCase(s.config, booksRepo, s.logger)

	// Init handlers
	booksHandlers := booksHttp.NewBooksHandlers(s.config, booksUC, s.logger)

	booksHttp.RegisterBookRoutes(router, booksHandlers)

	return router
}
