package tokens

import "bookloop.net/internal/models"

type Repository interface {
	Insert(token *models.Token) error
	DeleteAllForUser(scope string, userID int64) error
}
