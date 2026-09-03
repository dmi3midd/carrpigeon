package domain

import "time"

type Email struct {
	ID             string    `json:"id" db:"id"`
	Sender         string    `json:"sender" db:"sender"`
	Receiver       string    `json:"receiver" db:"receiver"`
	Subject        string    `json:"subject" db:"subject"`
	Body           string    `json:"body" db:"body"`
	IsHTML         bool      `json:"is_html" db:"is_html"`
	HTMLTemplateID *string   `json:"html_template_id" db:"html_template_id"`
	SentAt         time.Time `json:"sent_at" db:"sent_at"`
}
