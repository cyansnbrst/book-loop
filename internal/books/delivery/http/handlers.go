package http

import (
	"log/slog"
	"net/http"

	"bookloop.net/config"
	"bookloop.net/internal/books"
	"bookloop.net/internal/models"
	erp "bookloop.net/pkg/error_responses"
	"bookloop.net/pkg/utils"
	u "bookloop.net/pkg/utils"
	"bookloop.net/pkg/validator"
)

type booksHandlers struct {
	cfg     *config.Config
	booksUC books.UseCase
	logger  slog.Logger
}

func (h booksHandlers) listBooksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Title  string
			Author string
			Genres []string
			utils.Filters
		}

		v := validator.New()

		qs := r.URL.Query()

		input.Title = u.ReadString(qs, "title", "")
		input.Author = u.ReadString(qs, "author", "")
		input.Genres = u.ReadCSV(qs, "genres", []string{})

		input.Filters.Page = u.ReadInt(qs, "page", 1, v)
		input.Filters.PageSize = u.ReadInt(qs, "page_size", 20, v)

		input.Filters.Sort = u.ReadString(qs, "sort", "-created_at")
		input.Filters.SortSafelist = []string{"id", "title", "author", "created_at", "-id", "-title", "-author", "-created_at"}

		if u.ValidateFilters(v, input.Filters); !v.Valid() {
			erp.FailedValidationResponse(w, r, &h.logger, v.Errors)
			return
		}

		books, metadata, err := models.Book.GetAll(input.Title, input.Author, input.Genres, input.Filters)
		if err != nil {
			erp.ServerErrorResponse(w, r, &h.logger, err)
			return
		}

		err = u.WriteJSON(w, http.StatusOK, u.Envelope{"books": books, "metadata": metadata}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, &h.logger, err)
		}
	}
}
