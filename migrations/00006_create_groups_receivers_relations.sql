-- +goose Up
CREATE TABLE groups_receivers (
	group_id VARCHAR(20) NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
	receiver_id VARCHAR(20) NOT NULL REFERENCES email_receivers(id) ON DELETE CASCADE,
	PRIMARY KEY (group_id, receiver_id)
);

-- +goose Down
DROP TABLE groups_receivers;
