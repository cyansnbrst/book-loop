package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"bookloop.net/config"
	"bookloop.net/internal/books"
	"bookloop.net/internal/models"
	erp "bookloop.net/pkg/error_responses"
	u "bookloop.net/pkg/utils"
	"bookloop.net/pkg/validator"
)

type booksHandlers struct {
	cfg     *config.Config
	booksUC books.UseCase
	logger  slog.Logger
}

func NewBooksHandlers(cfg *config.Config, booksUC books.UseCase, logger slog.Logger) books.Handlers {
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
			switch {
			case errors.Is(err, validator.ErrJSONIsNotValid):
				erp.FailedValidationResponse(w, r, &h.logger, v.Errors)
			default:
				erp.ServerErrorResponse(w, r, &h.logger, err)
			}
			return
		}

		err = u.WriteJSON(w, http.StatusOK, u.Envelope{"books": books, "metadata": metadata}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, &h.logger, err)
		}
	}
}

func (h booksHandlers) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.CreateBook

		err := u.ReadJSON(w, r, &input)
		if err != nil {
			erp.BadRequestResponse(w, r, &h.logger, err)
			return
		}

		v := validator.New()

		book, err := h.booksUC.Insert(input, v)
		if err != nil {
			switch {
			case errors.Is(err, validator.ErrJSONIsNotValid):
				erp.FailedValidationResponse(w, r, &h.logger, v.Errors)
			default:
				erp.ServerErrorResponse(w, r, &h.logger, err)
			}
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

		err = u.WriteJSON(w, http.StatusCreated, u.Envelope{"book": book}, headers)
		if err != nil {
			erp.ServerErrorResponse(w, r, &h.logger, err)
		}
	}
}
