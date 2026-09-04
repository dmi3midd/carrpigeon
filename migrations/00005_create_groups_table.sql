-- +goose Up
CREATE TABLE groups (
	id VARCHAR(20) PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	description TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE groups;