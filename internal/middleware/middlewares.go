package middleware

import (
	"log/slog"

	"bookloop.net/config"
)

type MiddlewareManager struct {
	cfg    *config.Config
	logger *slog.Logger
}

func NewMiddlewareManager(cfg *config.Config, logger *slog.Logger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, logger: logger}
}
