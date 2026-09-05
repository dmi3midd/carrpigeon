package domain

import "time"

type EmailStatus string

const (
	StatusPending    EmailStatus = "pending"
	StatusProcessing EmailStatus = "processing"
	StatusSent       EmailStatus = "sent"
	StatusFailed     EmailStatus = "failed"
)

type Email struct {
	ID          string      `json:"id" db:"id"`
	Sender      string      `json:"sender" db:"sender"`
	Receiver    string      `json:"receiver" db:"receiver"`
	Subject     string      `json:"subject" db:"subject"`
	Body        string      `json:"body" db:"body"`
	TemplateID  *string     `json:"template_id" db:"template_id"`
	Status      EmailStatus `json:"status" db:"status"`
	Attempts    int         `json:"attempts" db:"attempts"`
	NextRetryAt *time.Time  `json:"next_retry_at" db:"next_retry_at"`
	LastError   *string     `json:"last_error" db:"last_error"`
	SentAt      *time.Time  `json:"sent_at" db:"sent_at"`
}
