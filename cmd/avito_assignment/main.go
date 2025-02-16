package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koccyx/avito_assignment/internal/config"
	"github.com/koccyx/avito_assignment/internal/server"
)


func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	log := setupLogger()

	log.Info("main started")
	log.Debug("debug messages enabled")
	
	serv := server.NewServer(log, cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	serv.SetupServer()
	
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serv.GracefulShutdown(ctx)
	
	log.Info("server stopped")
}


func setupLogger() *slog.Logger {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	return log
}