package postgres

import (
	"context"
	"database/sql"
	"time"

	"bookloop.net/config"
	_ "github.com/lib/pq"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PostgreSQL.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.PostgreSQL.MaxOpenConns)
	db.SetMaxIdleConns(cfg.PostgreSQL.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.PostgreSQL.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
