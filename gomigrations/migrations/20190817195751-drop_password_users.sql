-- +migrate Up
ALTER TABLE "users"
    DROP "password";
-- +migrate Down
