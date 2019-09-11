-- +migrate Up
ALTER TABLE refresh_tokens
    ADD COLUMN user_agent text not null;
-- +migrate Down
