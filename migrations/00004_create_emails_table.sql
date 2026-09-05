-- +goose Up
CREATE TABLE emails (
    id VARCHAR(20) PRIMARY KEY,
    sender VARCHAR(255) NOT NULL,
    receiver VARCHAR(255) NOT NULL REFERENCES receivers(email),
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    template_id VARCHAR(20) REFERENCES templates(id) ON DELETE SET NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS emails;