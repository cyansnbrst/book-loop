package error_responses

import (
	"fmt"
	"log/slog"
	"net/http"

	"bookloop.net/pkg/sl"
	"bookloop.net/pkg/utils"
)

func logError(r *http.Request, l *slog.Logger, err error) {
	l.Error("an error occured",
		slog.String("request_method", r.Method),
		slog.String("request_url", r.URL.String()),
		sl.Err(err),
	)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}, l *slog.Logger) {
	env := utils.Envelope{"error": message}

	err := utils.WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, l, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger, err error) {
	logError(r, l, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message, l)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message, l)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message, l)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error(), l)
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger, errors map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors, l)
}

func EditConflictResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	message := "unable to update the record due to an edit conflict, please try again"
	errorResponse(w, r, http.StatusConflict, message, l)
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	message := "rate limit exceeded"
	errorResponse(w, r, http.StatusTooManyRequests, message, l)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	message := "invalid authentication credentials"
	errorResponse(w, r, http.StatusUnauthorized, message, l)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message, l)
}

func AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	message := "you must be authenticated to access this resource"
	errorResponse(w, r, http.StatusUnauthorized, message, l)
}

func InactiveAccountResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	message := "your user account must be activated to access this resource"
	errorResponse(w, r, http.StatusForbidden, message, l)
}

func NotPermittedResponse(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	message := "your user account doesnt't have the necessare permissions to access this resource"
	errorResponse(w, r, http.StatusForbidden, message, l)
}
