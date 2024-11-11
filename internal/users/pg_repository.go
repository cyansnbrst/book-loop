package users

import "bookloop.net/internal/models"

type Repository interface {
	Insert(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	GetForToken(tokenScope, tokenPlaintext string) (*models.User, error)
}
