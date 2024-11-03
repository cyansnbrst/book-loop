package usecase

import (
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
	logger    *slog.Logger
}

func NewBooksUseCase(cfg *config.Config, booksRepo books.Repository, logger *slog.Logger) books.UseCase {
	return &booksUC{cfg: cfg, booksRepo: booksRepo, logger: logger}
}

func (u *booksUC) List(input models.BooksList, v *validator.Validator) ([]*models.Book, utils.Pagination, error) {
	input.Filters.SortSafelist = []string{"id", "title", "author", "created_at", "-id", "-title", "-author", "-created_at"}
	utils.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		return nil, utils.Pagination{}, validator.ErrJSONIsNotValid
	}

	books, metadata, err := u.booksRepo.GetAll(input.Title, input.Author, input.Genres, input.Filters)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	return books, metadata, nil
}

func (u *booksUC) Insert(input models.CreateBook, v *validator.Validator) (*models.Book, error) {
	book := &models.Book{
		Title:  input.Title,
		Author: input.Author,
		Genres: input.Genres,
	}

	if models.ValidateBook(v, book); !v.Valid() {
		return nil, validator.ErrJSONIsNotValid
	}

	err := u.booksRepo.Insert(book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (u *booksUC) Get(id int64) (*models.Book, error) {
	book, err := u.booksRepo.Get(id)
	return book, err
}
