package middleware

import (
	"fmt"
	"net/http"

	erp "bookloop.net/pkg/error_responses"
)

func (mw *MiddlewareManager) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				erp.ServerErrorResponse(w, r, mw.logger, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
