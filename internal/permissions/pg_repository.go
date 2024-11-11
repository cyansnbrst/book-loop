package permissions

import "bookloop.net/internal/models"

type Repository interface {
	GetAllForUser(userID int64) (models.Permissions, error)
	AddForUser(userID int64, codes ...string) error
}
