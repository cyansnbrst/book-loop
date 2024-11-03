package http

import (
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
		var list models.BooksList

		v := validator.New()

		qs := r.URL.Query()

		list.Title = u.ReadString(qs, "title", "")
		list.Author = u.ReadString(qs, "author", "")
		list.Genres = u.ReadCSV(qs, "genres", []string{})

		list.Filters.Page = u.ReadInt(qs, "page", 1, v)
		list.Filters.PageSize = u.ReadInt(qs, "page_size", 20, v)

		list.Filters.Sort = u.ReadString(qs, "sort", "-created_at")

		books, metadata, err := h.booksUC.ListBooks(list, v)
		if err != nil {
			if v.Errors != nil {
				erp.FailedValidationResponse(w, r, &h.logger, v.Errors)
			} else {
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
