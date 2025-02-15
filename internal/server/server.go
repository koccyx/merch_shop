package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/koccyx/avito_assignment/internal/config"
	"github.com/koccyx/avito_assignment/internal/server/handlers"
	"github.com/koccyx/avito_assignment/internal/server/middleware/auth"
	"github.com/koccyx/avito_assignment/internal/server/middleware/logger"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
	"github.com/koccyx/avito_assignment/internal/service"
	"github.com/koccyx/avito_assignment/internal/storage"
	"github.com/koccyx/avito_assignment/internal/storage/postgres"
)

type Server struct {
	Server *http.Server
	db *sql.DB
	log *slog.Logger
	cfg *config.Config
}

func (s *Server)SetupServer() {
	var err error
	s.db, err = postgres.New(s.cfg)

	if err != nil {
		s.log.Error("failed to init storage",sl.Err(err))
		os.Exit(1)
	}
	
	storage := storage.NewRepository(s.db)
	service := service.New(storage, s.log, s.cfg.Auth.Secret)

	router := chi.NewRouter()

	router.Use(logger.New(s.log))
	router.Use(middleware.Recoverer)
	
	router.Group(func(r chi.Router) {
		r.Post("/api/auth", handlers.Auth(service.Auth, s.log))
	})

	router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(service.Auth, s.cfg.Auth.Secret, s.log))
		r.Post("/api/sendCoin", handlers.TransferCoins(service.User, s.log))
		r.Get("/api/buy/{item}", handlers.PurchaseItem(service.Item, s.log))
		r.Get("/api/info", handlers.Info(service.User, s.log))
	})

	s.Server = &http.Server{
		Addr: fmt.Sprintf("%s:%s", s.cfg.Server.Addres, s.cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := s.Server.ListenAndServe(); err != nil {
			s.log.Error("failed to start server", sl.Err(err))
		}
	}()

	s.log.Info("server started")
}

func (s *Server) GracefulShutdown(ctx context.Context) {
	s.log.Info("graceful shutdown")

	if err := s.Server.Shutdown(ctx); err != nil {
		s.log.Error("failed to stop server", sl.Err(err))
		return
	}

	if err := s.db.Close(); err != nil {
		s.log.Error("failed to close db", sl.Err(err))
		return
	}
}

func NewServer(log *slog.Logger, cfg *config.Config) *Server{
	return &Server{
		log: log,
		cfg: cfg,
	}
}