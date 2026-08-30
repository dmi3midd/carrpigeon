package main

import (
	"carrpigeo/internal/config"
	"carrpigeo/internal/logs"
	"carrpigeo/internal/postgres"
	"carrpigeo/internal/repository"
	"carrpigeo/internal/server"
	"carrpigeo/internal/service"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"html/template"

	"github.com/dmi3midd/shkvcache"
)

func main() {
	// Root context with signal cancellation for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logFile, err := logs.Setup(cfg.Log.LogPath)
	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}
	defer logFile.Close()

	db, err := postgres.New(&cfg.Database)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Cache
	cache, err := shkvcache.NewCache[*template.Template](ctx, &shkvcache.Options{})
	if err != nil {
		slog.Error("failed to initialize cache", "error", err)
		os.Exit(1)
	}
	defer cache.Close()

	// Repositories
	htmlTemplateRepository := repository.NewHTMLTemplateRepository(db.GetDB())
	emailRepository := repository.NewEmailRepository(db.GetDB())

	// Services
	htmlTemplateService := service.NewHTMLTemplateService(htmlTemplateRepository, cache)
	emailClient := service.NewEmailClient(&cfg.SMTP)
	emailService := service.NewEmailService(emailClient, emailRepository, htmlTemplateService, &cfg.SMTP)

	server := server.NewServer(cfg, db, emailService, htmlTemplateService)

	slog.Info(
		"server is running",
		slog.String("address", cfg.HTTPServer.Address),
	)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error(
			"failed to run server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Graceful shutdown
	<-ctx.Done()
	slog.Info("received shutdown signal, stopping application...")

	server.Close()
	slog.Info("application stopped gracefully")
}
