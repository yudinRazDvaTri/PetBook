
-- +migrate Up
ALTER TABLE messages 
DROP COLUMN IF EXISTS read;
-- +migrate Down
ALTER TABLE messages
ADD COLUMN IF NOT EXISTS read boolean not null;