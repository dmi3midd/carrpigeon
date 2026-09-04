-- +goose Up
CREATE TYPE email_status AS ENUM ('pending', 'processing', 'sent', 'failed');

ALTER TABLE emails 
    ADD COLUMN status email_status NOT NULL DEFAULT 'pending',
    ADD COLUMN attempts INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN next_retry_at TIMESTAMPTZ,
    ADD COLUMN last_error TEXT,
    ALTER COLUMN sent_at DROP NOT NULL,
    ALTER COLUMN sent_at DROP DEFAULT;

CREATE INDEX idx_emails_status_next_retry ON emails (status, next_retry_at);

-- +goose Down
DROP INDEX IF EXISTS idx_emails_status_next_retry;
ALTER TABLE emails 
    DROP COLUMN IF EXISTS last_error,
    DROP COLUMN IF EXISTS next_retry_at,
    DROP COLUMN IF EXISTS attempts,
    DROP COLUMN IF EXISTS status,
    ALTER COLUMN sent_at SET DEFAULT CURRENT_TIMESTAMP;
DROP TYPE IF EXISTS email_status;


