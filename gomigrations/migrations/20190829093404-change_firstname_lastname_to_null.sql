
-- +migrate Up
ALTER TABLE users
    ALTER COLUMN firstname DROP NOT NULL,
    ALTER COLUMN lastname DROP NOT NULL;
-- +migrate Down
