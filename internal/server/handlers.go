package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	booksHttp "bookloop.net/internal/books/delivery/http"
	usersHttp "bookloop.net/internal/users/delivery/http"

	booksUseCase "bookloop.net/internal/books/usecase"
	permissionsUseCase "bookloop.net/internal/permissions/usecase"
	tokensUseCase "bookloop.net/internal/tokens/usecase"
	usersUseCase "bookloop.net/internal/users/usecase"

	booksRepository "bookloop.net/internal/books/repository"
	permissionsRepository "bookloop.net/internal/permissions/repository"
	tokensRepository "bookloop.net/internal/tokens/repository"
	usersRepository "bookloop.net/internal/users/repository"

	"bookloop.net/internal/middleware"
	"bookloop.net/pkg/error_responses"
)

func (s *Server) RegisterHandlers() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error_responses.NotFoundResponse(w, r, s.logger)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error_responses.MethodNotAllowedResponse(w, r, s.logger)
	})

	// Init repositories
	booksRepo := booksRepository.NewBooksRepository(s.db)
	usersRepo := usersRepository.NewUsersRepository(s.db)
	permissionsRepo := permissionsRepository.NewPermissionsRepository(s.db)
	tokensRepo := tokensRepository.NewTokensRepository(s.db)

	// Init useCases
	booksUC := booksUseCase.NewBooksUseCase(s.config, booksRepo, s.logger)
	usersUC := usersUseCase.NewUsersUseCase(s.config, usersRepo, s.logger)
	permissionsUC := permissionsUseCase.NewPermissionsUseCase(s.config, permissionsRepo, s.logger)
	tokensUC := tokensUseCase.NewTokensUseCase(s.config, tokensRepo, s.logger)

	// Init handlers
	booksHandlers := booksHttp.NewBooksHandlers(s.config, booksUC, s.logger)
	usersHandlers := usersHttp.NewUsersHandlers(s.config, usersUC, permissionsUC, tokensUC, &s.wg, s.logger, &s.mailer)

	mw := middleware.NewMiddlewareManager(s.config, usersUC, permissionsUC, s.logger)

	booksHttp.RegisterBookRoutes(router, booksHandlers, mw)
	usersHttp.RegisterUserRoutes(router, usersHandlers, mw)

	return mw.RecoverPanic(mw.RateLimit(mw.Authenticate(router)))
}
