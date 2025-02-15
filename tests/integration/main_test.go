package integration

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/koccyx/avito_assignment/internal/config"
	"github.com/koccyx/avito_assignment/internal/server"
)

var apiURL string
var cfg *config.Config

func TestMain(m *testing.M) {
	var err error

	cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}

	apiURL = "http://" + net.JoinHostPort(cfg.Server.Addres, cfg.Server.Port)

	log := setupLogger()

	log.Info("main started")
	log.Debug("debug messages enabled")
	
	serv := server.NewServer(log, cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	serv.SetupServer()
	
	exitCode := m.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serv.GracefulShutdown(ctx)

    println("Tearing down...")

    os.Exit(exitCode)
}

func setupLogger() *slog.Logger {
	log := slog.New(
		slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	
	return log
}