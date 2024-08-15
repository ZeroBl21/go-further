package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("shutting down server", slog.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks", slog.String("addr", srv.Addr))

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("starting server",
		slog.String("addr", srv.Addr),
		slog.String("env", app.config.env),
	)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err := <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", slog.String("addr", srv.Addr))

	return nil
}
