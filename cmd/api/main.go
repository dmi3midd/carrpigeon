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
	htmltemplate "html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	txttemplate "text/template"

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

	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Cache
	parsedHtmlTmplCache, err := shkvcache.NewCache[*htmltemplate.Template](ctx, &shkvcache.Options{
		ShardCount:      4,
		CleanerInterval: 120,
		RunCleaner:      true,
	})
	if err != nil {
		slog.Error("failed to initialize parsed html template cache", "error", err)
		os.Exit(1)
	}
	defer parsedHtmlTmplCache.Close()
	parsedTxtTmplCache, err := shkvcache.NewCache[*txttemplate.Template](ctx, &shkvcache.Options{
		ShardCount:      4,
		CleanerInterval: 120,
		RunCleaner:      true,
	})
	if err != nil {
		slog.Error("failed to initialize parsed txt template cache", "error", err)
		os.Exit(1)
	}
	defer parsedTxtTmplCache.Close()
	domainTmplCache, err := shkvcache.NewCache[*domain.Template](ctx, &shkvcache.Options{
		ShardCount:      4,
		CleanerInterval: 120,
		RunCleaner:      true,
	})
	if err != nil {
		slog.Error("failed to initialize raw template cache", "error", err)
		os.Exit(1)
	}
	defer domainTmplCache.Close()

	emailCache, err := shkvcache.NewCache[*domain.Email](ctx, &shkvcache.Options{
		ShardCount:      8,
		CleanerInterval: 60,
		RunCleaner:      false,
	})
	if err != nil {
		slog.Error("failed to initialize email cache", "error", err)
		os.Exit(1)
	}
	defer emailCache.Close()

	// Repositories
	templateRepository := repository.NewTemplateRepository(db.GetDB())
	emailRepository := repository.NewEmailRepository(db.GetDB())
	emailReceiverRepository := repository.NewReceiverRepository(db.GetDB())
	groupRepository := repository.NewGroupRepository(db.GetDB())

	templateService := service.NewTemplateService(templateRepository, parsedHtmlTmplCache, parsedTxtTmplCache, domainTmplCache)
	emailClient := client.NewEmailClient(&cfg.Email.SMTP)
	emailService := service.NewEmailService(emailClient, emailRepository, emailReceiverRepository, templateService, &cfg.Email.SMTP)
	emailReceiverService := service.NewReceiverService(emailReceiverRepository)
	groupService := service.NewGroupService(groupRepository, emailReceiverRepository)

	// Worker
	emailWorker := client.NewEmailWorker(emailClient, emailRepository, emailCache, cfg.Email.Worker)
	emailWorker.Start(ctx)
	defer emailWorker.Stop()

	// Handlers
	emailReceiversHandler := handlers.NewReceiversHandler(emailReceiverService)
	sendHandler := handlers.NewSendHandler(emailService)
	templateHandler := handlers.NewTemplateHandler(templateService)
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
	emailWorker.Stop()
	slog.Info("application stopped gracefully")
}
