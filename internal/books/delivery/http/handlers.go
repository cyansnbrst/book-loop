package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"bookloop.net/config"
	"bookloop.net/internal/books"
	"bookloop.net/internal/models"
	"bookloop.net/pkg/db"
	erp "bookloop.net/pkg/error_responses"
	u "bookloop.net/pkg/utils"
	"bookloop.net/pkg/validator"
)

type booksHandlers struct {
	cfg     *config.Config
	booksUC books.UseCase
	logger  *slog.Logger
}

func NewBooksHandlers(cfg *config.Config, booksUC books.UseCase, logger *slog.Logger) books.Handlers {
	return &booksHandlers{cfg: cfg, booksUC: booksUC, logger: logger}
}

func (h booksHandlers) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.BooksList

		v := validator.New()

		qs := r.URL.Query()

		input.Title = u.ReadString(qs, "title", "")
		input.Author = u.ReadString(qs, "author", "")
		input.Genres = u.ReadCSV(qs, "genres", []string{})

		input.Filters.Page = u.ReadInt(qs, "page", 1, v)
		input.Filters.PageSize = u.ReadInt(qs, "page_size", 20, v)

		input.Filters.Sort = u.ReadString(qs, "sort", "-created_at")

		books, metadata, err := h.booksUC.List(input, v)
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

		err = u.WriteJSON(w, http.StatusOK, u.Envelope{"books": books, "metadata": metadata}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

func (h booksHandlers) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.InputBook

		err := u.ReadJSON(w, r, &input)
		if err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		book, err := h.booksUC.Insert(input)
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

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

		err = u.WriteJSON(w, http.StatusCreated, u.Envelope{"book": book}, headers)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

func (h booksHandlers) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := u.ReadIDParam(r)
		if err != nil {
			erp.NotFoundResponse(w, r, h.logger)
			return
		}

		book, err := h.booksUC.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrRecordNotFound):
				erp.NotFoundResponse(w, r, h.logger)
			default:
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		err = u.WriteJSON(w, http.StatusOK, u.Envelope{"book": book}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

func (h booksHandlers) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := u.ReadIDParam(r)
		if err != nil {
			erp.NotFoundResponse(w, r, h.logger)
			return
		}

		var input models.InputBook

		err = u.ReadJSON(w, r, &input)
		if err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		book, err := h.booksUC.Update(id, input)
		if err != nil {
			var validationErr *validator.ValidationError
			switch {
			case errors.As(err, &validationErr):
				erp.FailedValidationResponse(w, r, h.logger, validationErr.Errors)
			case errors.Is(err, db.ErrRecordNotFound):
				erp.NotFoundResponse(w, r, h.logger)
			case errors.Is(err, db.ErrEditConflict):
				erp.EditConflictResponse(w, r, h.logger)
			default:
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		err = u.WriteJSON(w, http.StatusOK, u.Envelope{"book": book}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}

	}
}
