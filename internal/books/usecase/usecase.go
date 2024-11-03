package usecase

import (
	"fmt"
	"log/slog"

	"bookloop.net/config"
	"bookloop.net/internal/books"
	"bookloop.net/internal/models"
	"bookloop.net/pkg/utils"
	"bookloop.net/pkg/validator"
)

type booksUC struct {
	cfg       *config.Config
	booksRepo books.Repository
	logger    slog.Logger
}

func NewBooksUseCase(cfg *config.Config, booksRepo books.Repository, logger slog.Logger) books.UseCase {
	return &booksUC{cfg: cfg, booksRepo: booksRepo, logger: logger}
}

func (u *booksUC) ListBooks(list models.BooksList, v *validator.Validator) ([]*models.Book, utils.Pagination, error) {
	list.Filters.SortSafelist = []string{"id", "title", "author", "created_at", "-id", "-title", "-author", "-created_at"}
	utils.ValidateFilters(v, list.Filters)

	if !v.Valid() {
		return nil, utils.Pagination{}, fmt.Errorf("validation error: %v", v.Errors)
	}

	books, metadata, err := u.booksRepo.GetAll(list.Title, list.Author, list.Genres, list.Filters)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	return books, metadata, nil
}
