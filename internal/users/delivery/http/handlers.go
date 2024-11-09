package http

import (
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"bookloop.net/config"
	"bookloop.net/internal/models"
	"bookloop.net/internal/permissions"
	"bookloop.net/internal/tokens"
	"bookloop.net/internal/users"
	"bookloop.net/pkg/db"
	erp "bookloop.net/pkg/error_responses"
	"bookloop.net/pkg/mailer"
	"bookloop.net/pkg/sl"
	u "bookloop.net/pkg/utils"
	"bookloop.net/pkg/validator"
)

type usersHandlers struct {
	cfg           *config.Config
	usersUC       users.UseCase
	permissionsUC permissions.UseCase
	tokensUC      tokens.UseCase
	wg            *sync.WaitGroup
	logger        *slog.Logger
	mailer        *mailer.Mailer
}

func NewUsersHandlers(cfg *config.Config, usersUC users.UseCase, permissionsUC permissions.UseCase, tokensUC tokens.UseCase, wg *sync.WaitGroup, logger *slog.Logger, mailer *mailer.Mailer) users.Handlers {
	return &usersHandlers{cfg: cfg, usersUC: usersUC, permissionsUC: permissionsUC, tokensUC: tokensUC, wg: wg, logger: logger, mailer: mailer}
}

func (h usersHandlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.InputRegisterUser

		err := u.ReadJSON(w, r, &input)
		if err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		user, err := h.usersUC.Register(input)
		if err != nil {
			var validationErr *validator.ValidationError
			switch {
			case errors.As(err, &validationErr):
				erp.FailedValidationResponse(w, r, h.logger, validationErr.Errors)
			default:
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		err = h.permissionsUC.AddForUser(user.ID, "books:read")
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		token, err := h.tokensUC.NewActivationToken(user.ID)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		u.Background(h.wg, h.logger, func() {
			data := map[string]interface{}{
				"activationToken": token.Plaintext,
				"userID":          user.ID,
			}

			err = h.mailer.Send(user.Email, "user_welcome.tmpl", data)
			if err != nil {
				h.logger.Error("an error occured while sending email", sl.Err(err))
			}
		})

		err = u.WriteJSON(w, http.StatusAccepted, u.Envelope{"user": user}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

func (h usersHandlers) Activate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.InputToken

		err := u.ReadJSON(w, r, &input)
		if err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		user, err := h.usersUC.Activate(input.TokenPlaintext)
		if err != nil {
			var validationErr *validator.ValidationError
			switch {
			case errors.As(err, &validationErr):
				erp.FailedValidationResponse(w, r, h.logger, validationErr.Errors)
			case errors.Is(err, db.ErrEditConflict):
				erp.EditConflictResponse(w, r, h.logger)
			default:
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		err = h.tokensUC.DeleteActivationTokens(user.ID)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		err = u.WriteJSON(w, http.StatusOK, u.Envelope{"user": user}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

func (h usersHandlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.InputLoginUser

		err := u.ReadJSON(w, r, &input)
		if err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		user, err := h.usersUC.Login(input)
		if err != nil {
			var validationErr *validator.ValidationError
			switch {
			case errors.As(err, &validationErr):
				erp.FailedValidationResponse(w, r, h.logger, validationErr.Errors)
			case errors.Is(err, db.ErrRecordNotFound):
				erp.InvalidCredentialsResponse(w, r, h.logger)
			default:
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		token, err := h.tokensUC.NewAuthenticationToken(user.ID)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		err = u.WriteJSON(w, http.StatusCreated, u.Envelope{"authentication_token": token}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

func (h usersHandlers) CreateActivationToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.InputNewActivationTokenUser

		err := u.ReadJSON(w, r, &input)
		if err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		user, err := h.usersUC.NewActivationToken(input)
		if err != nil {
			var validationErr *validator.ValidationError
			switch {
			case errors.As(err, &validationErr):
				erp.FailedValidationResponse(w, r, h.logger, validationErr.Errors)
			default:
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		token, err := h.tokensUC.NewActivationToken(user.ID)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		u.Background(h.wg, h.logger, func() {
			data := map[string]interface{}{
				"activationToken": token.Plaintext,
			}

			err = h.mailer.Send(user.Email, "token_activation.tmpl", data)
			if err != nil {
				h.logger.Error("an error occured while sending email", sl.Err(err))
			}
		})

		env := u.Envelope{"message": "an email will be sent to you containing activation instructions"}

		err = u.WriteJSON(w, http.StatusAccepted, env, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}
