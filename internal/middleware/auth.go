package middleware

import (
	"errors"
	"net/http"
	"strings"

	"bookloop.net/internal/models"
	"bookloop.net/pkg/db"
	erp "bookloop.net/pkg/error_responses"
	"bookloop.net/pkg/validator"
)

func (mw *MiddlewareManager) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = ContextSetUser(r, models.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			erp.InvalidAuthenticationTokenResponse(w, r, mw.logger)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if models.ValidateTokenPlainText(v, token); !v.Valid() {
			erp.InvalidAuthenticationTokenResponse(w, r, mw.logger)
			return
		}

		user, err := mw.usersUC.Authenticate(token)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrRecordNotFound):
				erp.InvalidAuthenticationTokenResponse(w, r, mw.logger)
			default:
				erp.ServerErrorResponse(w, r, mw.logger, err)
			}
			return
		}

		r = ContextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (mw *MiddlewareManager) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ContextGetUser(r)
		if user.IsAnonymous() {
			erp.AuthenticationRequiredResponse(w, r, mw.logger)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (mw *MiddlewareManager) RequireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ContextGetUser(r)
		if !user.Activated {
			erp.InactiveAccountResponse(w, r, mw.logger)
			return
		}
		next.ServeHTTP(w, r)
	})

	return mw.RequireAuthenticatedUser(fn)
}

func (mw *MiddlewareManager) RequirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := ContextGetUser(r)

		permissions, err := mw.permissionsUC.GetAllForUser(user.ID)
		if err != nil {
			erp.ServerErrorResponse(w, r, mw.logger, err)
			return
		}

		if !permissions.Include(code) {
			erp.NotPermittedResponse(w, r, mw.logger)
			return
		}
		next.ServeHTTP(w, r)
	}

	return mw.RequireActivatedUser(fn)
}
