package usecase

import (
	"log/slog"

	"bookloop.net/config"
	"bookloop.net/internal/models"
	"bookloop.net/internal/permissions"
)

type permissionsUC struct {
	cfg             *config.Config
	permissionsRepo permissions.Repository
	logger          *slog.Logger
}

func NewPermissionsUseCase(cfg *config.Config, permissionsRepo permissions.Repository, logger *slog.Logger) permissions.UseCase {
	return &permissionsUC{cfg: cfg, permissionsRepo: permissionsRepo, logger: logger}
}

func (u *permissionsUC) GetAllForUser(userID int64) (models.Permissions, error) {
	return u.permissionsRepo.GetAllForUser(userID)
}

func (u *permissionsUC) AddForUser(userID int64, codes ...string) error {
	return u.permissionsRepo.AddForUser(userID, codes...)
}
