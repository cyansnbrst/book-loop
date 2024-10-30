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
	Config config.Config
	Logger slog.Logger
	DB     *sql.DB
	Mailer mailer.Mailer
	WG     sync.WaitGroup
}

func (s *Server) Serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.Port),
		Handler:      s.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutDownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		s.Logger.Info("shutting down server",
			slog.String("signal", sig.String()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutDownError <- err
		}

		s.Logger.Info("completing background tasks",
			slog.String("addr", server.Addr),
		)

		s.WG.Wait()
		shutDownError <- nil
	}()

	s.Logger.Info("starting server",
		slog.String("addr", server.Addr),
		slog.String("env", s.Config.Env),
	)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownError
	if err != nil {
		return err
	}

	s.Logger.Info("stopped server",
		slog.String("addr", server.Addr),
	)

	return nil
}
