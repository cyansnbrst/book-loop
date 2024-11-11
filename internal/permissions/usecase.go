package permissions

import "bookloop.net/internal/models"

type UseCase interface {
	GetAllForUser(userID int64) (models.Permissions, error)
	AddForUser(userID int64, codes ...string) error
}
