package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"bookloop.net/config"
	"bookloop.net/pkg/mailer"
)

type Server struct {
	config *config.Config
	logger slog.Logger
	db     *sql.DB
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func NewServer(cfg *config.Config, logger slog.Logger, db *sql.DB, mailer mailer.Mailer) *Server {
	return &Server{config: cfg, logger: logger, db: db, mailer: mailer}
}

func (s *Server) Serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      s.RegisterHandlers(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutDownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		s.logger.Info("shutting down server",
			slog.String("signal", sig.String()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutDownError <- err
		}

		s.logger.Info("completing background tasks",
			slog.String("addr", server.Addr),
		)

		s.wg.Wait()
		shutDownError <- nil
	}()

	s.logger.Info("starting server",
		slog.String("addr", server.Addr),
		slog.String("env", s.config.Env),
	)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownError
	if err != nil {
		return err
	}

	s.logger.Info("stopped server",
		slog.String("addr", server.Addr),
	)

	return nil
}
