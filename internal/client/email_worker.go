package client

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/dmi3midd/carrpigeon/internal/config"
	"github.com/dmi3midd/carrpigeon/internal/domain"
	"github.com/dmi3midd/carrpigeon/internal/repository"

	"github.com/dmi3midd/shkvcache"
)

type EmailWorker interface {
	Start(ctx context.Context)
	Stop()
	RecoverStuckEmails(ctx context.Context) error
}

type emailWorker struct {
	client     EmailClient
	emailRepo  repository.EmailRepository
	emailCache *shkvcache.Cache[*domain.Email]
	cfg        config.Worker

	jobsChan chan domain.Email
	wg       sync.WaitGroup
	stopOnce sync.Once
	stopChan chan struct{}
}

func NewEmailWorker(
	client EmailClient,
	emailRepo repository.EmailRepository,
	emailCache *shkvcache.Cache[*domain.Email],
	cfg config.Worker,
) EmailWorker {
	if cfg.WorkerPoolSize <= 0 {
		cfg.WorkerPoolSize = 5
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 2 * time.Second
	}
	if cfg.RetryDelay <= 0 {
		cfg.RetryDelay = 10 * time.Second
	}

	return &emailWorker{
		client:     client,
		emailRepo:  emailRepo,
		emailCache: emailCache,
		cfg:        cfg,
		jobsChan:   make(chan domain.Email, cfg.WorkerPoolSize*2),
		stopChan:   make(chan struct{}),
	}
}

func (w *emailWorker) RecoverStuckEmails(ctx context.Context) error {
	const batchSize = 50
	totalRecovered := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		emails, err := w.emailRepo.ResetProcessing(ctx, batchSize)
		if err != nil {
			return fmt.Errorf("failed to recover stuck processing emails: %w", err)
		}

		if len(emails) == 0 {
			break
		}

		for _, email := range emails {
			totalRecovered++
			slog.Warn("recovered stuck email, reset status to pending",
				slog.String("id", email.ID),
				slog.String("receiver", email.ReceiverEmail),
			)
		}

		if len(emails) < batchSize {
			break
		}
	}

	if totalRecovered > 0 {
		slog.Info("stuck emails recovery completed", slog.Int("recovered_count", totalRecovered))
	}

	return nil
}

func (w *emailWorker) Start(ctx context.Context) {
	slog.Info("starting email worker pool",
		slog.Int("pool_size", w.cfg.WorkerPoolSize),
		slog.Int("max_retries", w.cfg.MaxRetries),
		slog.Duration("poll_interval", w.cfg.PollInterval),
	)

	if err := w.RecoverStuckEmails(ctx); err != nil {
		slog.Error("failed to recover stuck emails on startup", slog.String("error", err.Error()))
	}

	for i := 1; i <= w.cfg.WorkerPoolSize; i++ {
		w.wg.Add(1)
		go w.runWorker(ctx)
	}

	w.wg.Add(1)
	go w.runDispatcher(ctx)
}

func (w *emailWorker) runDispatcher(ctx context.Context) {
	defer w.wg.Done()
	defer close(w.jobsChan)

	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	w.dispatch(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopChan:
			return
		case <-ticker.C:
			w.dispatch(ctx)
		}
	}
}

func (w *emailWorker) dispatch(ctx context.Context) {
	batchSize := w.cfg.WorkerPoolSize * 2

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopChan:
			return
		default:
		}

		emails, err := w.emailRepo.FetchPending(ctx, batchSize)
		if err != nil {
			slog.Error("failed to fetch pending emails", slog.String("error", err.Error()))
			return
		}

		if len(emails) == 0 {
			return
		}

		for _, email := range emails {
			if w.emailCache != nil {
				_ = w.emailCache.Set(email.ID, &email, 300)
			}

			select {
			case <-ctx.Done():
				return
			case <-w.stopChan:
				return
			case w.jobsChan <- email:
			}
		}

		if len(emails) < batchSize {
			return
		}
	}
}

func (w *emailWorker) runWorker(ctx context.Context) {
	defer w.wg.Done()

	for email := range w.jobsChan {
		w.processEmail(ctx, email)
	}
}

func (w *emailWorker) processEmail(ctx context.Context, email domain.Email) {
	err := w.client.Send(&email)

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err == nil {
		now := time.Now()
		if updateErr := w.emailRepo.MarkAsSent(dbCtx, email.ID, now); updateErr != nil {
			slog.Error("failed to mark email as sent in db", slog.String("id", email.ID), slog.String("error", updateErr.Error()))
		} else {
			slog.Info("email sent successfully", slog.String("id", email.ID), slog.String("receiver", email.ReceiverEmail))
		}
		if w.emailCache != nil {
			w.emailCache.Del(email.ID)
		}
		return
	}

	slog.Error("failed to send email", slog.String("id", email.ID), slog.String("receiver", email.ReceiverEmail), slog.String("error", err.Error()))

	nextAttempts := email.Attempts + 1
	var nextRetryAt *time.Time

	if nextAttempts < w.cfg.MaxRetries {
		backoffDuration := w.cfg.RetryDelay * time.Duration(math.Pow(2, float64(email.Attempts)))
		nextTime := time.Now().Add(backoffDuration)
		nextRetryAt = &nextTime
		slog.Warn("scheduling email retry",
			slog.String("id", email.ID),
			slog.Int("attempt", nextAttempts),
			slog.Time("next_retry_at", nextTime),
		)
	} else {
		slog.Error("email reached max retry attempts, marking as failed", slog.String("id", email.ID), slog.Int("attempts", nextAttempts))
	}

	if updateErr := w.emailRepo.MarkAsFailed(dbCtx, email.ID, nextAttempts, nextRetryAt, err.Error()); updateErr != nil {
		slog.Error("failed to mark email status in db", slog.String("id", email.ID), slog.String("error", updateErr.Error()))
	}

	if w.emailCache != nil {
		w.emailCache.Del(email.ID)
	}
}

func (w *emailWorker) Stop() {
	w.stopOnce.Do(func() {
		slog.Info("stopping email worker...")
		close(w.stopChan)
		w.wg.Wait()
		if err := w.client.Close(); err != nil {
			slog.Error("failed to close email client", slog.String("error", err.Error()))
		}
		slog.Info("email worker stopped gracefully")
	})
}
