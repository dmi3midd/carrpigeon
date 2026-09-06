package main

import (
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

	"github.com/dmi3midd/carrpigeon/internal/client"
	"github.com/dmi3midd/carrpigeon/internal/config"
	"github.com/dmi3midd/carrpigeon/internal/domain"
	"github.com/dmi3midd/carrpigeon/internal/logger"
	"github.com/dmi3midd/carrpigeon/internal/postgres"
	"github.com/dmi3midd/carrpigeon/internal/repository"
	"github.com/dmi3midd/carrpigeon/internal/server"
	"github.com/dmi3midd/carrpigeon/internal/server/handlers"
	"github.com/dmi3midd/carrpigeon/internal/server/middlewares"
	"github.com/dmi3midd/carrpigeon/internal/service"

	"github.com/dmi3midd/shkvcache"
	"github.com/go-playground/validator/v10"
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
	receiverRepository := repository.NewReceiverRepository(db.GetDB())
	groupRepository := repository.NewGroupRepository(db.GetDB())
	txManager := repository.NewTxManager(db.GetDB())

	templateService := service.NewTemplateService(templateRepository, parsedHtmlTmplCache, parsedTxtTmplCache, domainTmplCache)
	sendService := service.NewSendService(emailRepository, receiverRepository, groupRepository, txManager, templateService, &cfg.Email.SMTP)
	receiverService := service.NewReceiverService(receiverRepository)
	groupService := service.NewGroupService(groupRepository, receiverRepository)

	// Worker & Client
	emailClient := client.NewEmailClient(&cfg.Email.SMTP)
	emailWorker := client.NewEmailWorker(emailClient, emailRepository, emailCache, cfg.Email.Worker)
	emailWorker.Start(ctx)
	defer emailWorker.Stop()

	validate := validator.New()

	// Handlers
	receiversHandler := handlers.NewReceiversHandler(receiverService, validate)
	sendHandler := handlers.NewSendHandler(sendService, validate)
	templateHandler := handlers.NewTemplateHandler(templateService)
	groupHandler := handlers.NewGroupHandler(groupService, validate)
	systemHandler := handlers.NewSystemHandler(db)

	// Middleware
	middlewares := middlewares.NewMiddlewares(cfg)

	server := server.NewServer(
		cfg,
		middlewares,
		systemHandler,
		sendHandler,
		receiversHandler,
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
