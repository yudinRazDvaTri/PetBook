
-- +migrate Up
ALTER TABLE topics
DROP comments_number;
-- +migrate Down
ALTER TABLE topics
ADD comments_number integer default 0;