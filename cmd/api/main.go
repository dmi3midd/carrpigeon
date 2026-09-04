package main

import (
	"carrpigeo/internal/client"
	"carrpigeo/internal/config"
	"carrpigeo/internal/domain"
	"carrpigeo/internal/logger"
	"carrpigeo/internal/postgres"
	"carrpigeo/internal/repository"
	"carrpigeo/internal/server"
	"carrpigeo/internal/server/handlers"
	"carrpigeo/internal/server/middlewares"
	"carrpigeo/internal/service"
	"context"
	"errors"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	// Initialize logger
	logger.Setup(cfg.Logger.Level)

	db, err := postgres.New(&cfg.Database)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Cache
	parsedTmplCache, err := shkvcache.NewCache[*template.Template](ctx, &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 60,
		RunCleaner:      true,
	})
	if err != nil {
		slog.Error("failed to initialize parsed template cache", "error", err)
		os.Exit(1)
	}
	defer parsedTmplCache.Close()

	rawTmplCache, err := shkvcache.NewCache[*domain.HTMLTemplate](ctx, &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 60,
		RunCleaner:      true,
	})
	if err != nil {
		slog.Error("failed to initialize raw template cache", "error", err)
		os.Exit(1)
	}
	defer rawTmplCache.Close()

	// Repositories
	htmlTemplateRepository := repository.NewHTMLTemplateRepository(db.GetDB())
	emailRepository := repository.NewEmailRepository(db.GetDB())
	emailReceiverRepository := repository.NewEmailReceiverRepository(db.GetDB())
	groupRepository := repository.NewGroupRepository(db.GetDB())

	// Services
	htmlTemplateService := service.NewHTMLTemplateService(htmlTemplateRepository, parsedTmplCache, rawTmplCache)
	emailClient := client.NewEmailClient(&cfg.SMTP)
	emailService := service.NewEmailService(emailClient, emailRepository, emailReceiverRepository, htmlTemplateService, &cfg.SMTP)
	emailReceiverService := service.NewEmailReceiverService(emailReceiverRepository)
	groupService := service.NewGroupService(groupRepository, emailReceiverRepository)

	// Handlers
	emailReceiversHandler := handlers.NewEmailReceiversHandler(emailReceiverService)
	sendHandler := handlers.NewSendHandler(emailService)
	templateHandler := handlers.NewTemplateHandlers(htmlTemplateService)
	groupHandler := handlers.NewGroupHandler(groupService)
	systemHandler := handlers.NewSystemHandler(db)

	// Middleware
	middlewares := middlewares.NewMiddlewares(cfg)

	server := server.NewServer(
		cfg,
		middlewares,
		systemHandler,
		sendHandler,
		emailReceiversHandler,
		templateHandler,
		groupHandler,
	)

	slog.Info(
		"server is running",
		slog.String("address", cfg.HTTPServer.Address),
	)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to run server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	slog.Info("received shutdown signal, stopping application...")

	server.Close()
	slog.Info("application stopped gracefully")
}
