-- +migrate Up
ALTER TABLE refresh_tokens
    ADD UNIQUE (user_id, user_agent);
-- +migrate Down
