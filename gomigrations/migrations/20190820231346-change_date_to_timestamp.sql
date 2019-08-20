
-- +migrate Up
ALTER TABLE messages 
	ALTER COLUMN created_at TYPE TIMESTAMP;

-- +migrate Down
ALTER TABLE messages 
	ALTER COLUMN created_at TYPE date;
