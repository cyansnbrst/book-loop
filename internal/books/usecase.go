package books

import (
	"bookloop.net/internal/models"
	"bookloop.net/pkg/utils"
	"bookloop.net/pkg/validator"
)

type UseCase interface {
	List(input models.BooksList, v *validator.Validator) ([]*models.Book, utils.Pagination, error)
	Insert(input models.InputBook) (*models.Book, error)
	Get(id int64) (*models.Book, error)
	Update(id int64, input models.InputBook) (*models.Book, error)
	Delete(id int64) error
}
