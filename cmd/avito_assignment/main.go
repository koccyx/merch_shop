package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/koccyx/avito_assignment/internal/config"
	"github.com/koccyx/avito_assignment/internal/http/handlers"
	"github.com/koccyx/avito_assignment/internal/http/middleware/auth"
	"github.com/koccyx/avito_assignment/internal/http/middleware/logger"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
	"github.com/koccyx/avito_assignment/internal/service"
	"github.com/koccyx/avito_assignment/internal/storage"
	"github.com/koccyx/avito_assignment/internal/storage/postgres"
)

const(
	envLocal = "local"
    envDev = "dev"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	log := setupLogger(cfg.Env)

	log.Info("main started")
	log.Debug("debug messages enabled")

	db, err := postgres.New(cfg)

	if err != nil {
		log.Error("failed to init storage",sl.Err(err))
		os.Exit(1)
	}
	
	storage := storage.NewRepository(db)
	service := service.New(storage, log, cfg.Auth.Secret)

	router := chi.NewRouter()

	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	
	router.Group(func(r chi.Router) {
		r.Post("/api/auth", handlers.Auth(service.Auth, log))
	})

	router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(service.Auth, cfg.Auth.Secret, log))
		r.Get("/api/buy/{item}", handlers.Merch(service.Auth, log))
	})

	serv := &http.Server{
		Addr: cfg.Server.Addres,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := serv.ListenAndServe(); err != nil {
			log.Error("failed to start server", sl.Err(err))
		}
	}()

	log.Info("server started")
	
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	db.Close()

	log.Info("server stopped")
}


func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}