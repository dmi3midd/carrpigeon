-- +goose Up
CREATE TYPE email_status AS ENUM ('pending', 'processing', 'sent', 'failed');
CREATE TABLE emails (
    id VARCHAR(20) PRIMARY KEY,
    sender VARCHAR(255) NOT NULL,
    receiver_id VARCHAR(20) REFERENCES receivers(id) ON DELETE SET NULL,
    receiver_email VARCHAR(255) NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    template_id VARCHAR(20) REFERENCES templates(id) ON DELETE SET NULL,
    status email_status NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMPTZ,
    last_error TEXT,
    sent_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_emails_status_next_retry ON emails (status, next_retry_at);

-- +goose Down
DROP INDEX IF EXISTS idx_emails_status_next_retry;
DROP TABLE IF EXISTS emails;
DROP TYPE IF EXISTS email_status;