package tokens

import (
	"time"

	"bookloop.net/internal/models"
)

type UseCase interface {
	New(userID int64, ttl time.Duration, scope string) (*models.Token, error)
	NewActivationToken(userID int64) (*models.Token, error)
	NewAuthenticationToken(userID int64) (*models.Token, error)
	DeleteActivationTokens(userID int64) error
}
