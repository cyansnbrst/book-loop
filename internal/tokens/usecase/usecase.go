package usecase

import (
	"log/slog"
	"time"

	"bookloop.net/config"
	"bookloop.net/internal/models"
	"bookloop.net/internal/tokens"
)

type tokensUC struct {
	cfg        *config.Config
	tokensRepo tokens.Repository
	logger     *slog.Logger
}

func NewTokensUseCase(cfg *config.Config, tokensRepo tokens.Repository, logger *slog.Logger) tokens.UseCase {
	return &tokensUC{cfg: cfg, tokensRepo: tokensRepo, logger: logger}
}

func (u *tokensUC) New(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	token, err := models.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = u.tokensRepo.Insert(token)
	return token, err
}

func (u *tokensUC) NewActivationToken(userID int64) (*models.Token, error) {
	return u.New(userID, 3*24*time.Hour, models.ScopeActivation)
}

func (u *tokensUC) NewAuthenticationToken(userID int64) (*models.Token, error) {
	return u.New(userID, 24*time.Hour, models.ScopeAuthentication)
}

func (u *tokensUC) DeleteActivationTokens(userID int64) error {
	return u.tokensRepo.DeleteAllForUser(models.ScopeActivation, userID)
}
