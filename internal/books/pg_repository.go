package books

import (
	"bookloop.net/internal/models"
	"bookloop.net/pkg/utils"
)

type Repository interface {
	GetAll(title string, author string, genres []string, filters utils.Filters) ([]*models.Book, utils.Pagination, error)
	Insert(book *models.Book) error
	Get(id int64) (*models.Book, error)
	Update(book *models.Book) error
	Delete(id int64) error
}
