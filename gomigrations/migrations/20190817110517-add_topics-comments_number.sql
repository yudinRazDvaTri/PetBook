
-- +migrate Up
ALTER TABLE topics
ADD comments_number integer default 0;
-- +migrate Down
ALTER TABLE topics
DROP comments_number;