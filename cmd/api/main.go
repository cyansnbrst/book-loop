package main

import (
	"log/slog"
	"os"

	"bookloop.net/config"
	"bookloop.net/internal/server"
	"bookloop.net/pkg/db/postgres"
	"bookloop.net/pkg/mailer"
	"bookloop.net/pkg/sl"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	db, err := postgres.OpenDB(cfg)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	s := *server.NewServer(
		cfg,
		logger,
		db,
		mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
	)

	err = s.Serve()
	if err != nil {
		logger.Error("an error occured", sl.Err(err))
		os.Exit(1)
	}
}
