package middleware

import (
	"log/slog"

	"bookloop.net/config"
	"bookloop.net/internal/permissions"
	"bookloop.net/internal/users"
)

type MiddlewareManager struct {
	cfg           *config.Config
	usersUC       users.UseCase
	permissionsUC permissions.UseCase
	logger        *slog.Logger
}

func NewMiddlewareManager(cfg *config.Config, usersUC users.UseCase, permissionsUC permissions.UseCase, logger *slog.Logger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, usersUC: usersUC, permissionsUC: permissionsUC, logger: logger}
}
