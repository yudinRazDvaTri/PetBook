
-- +migrate Up
ALTER TABLE comments
ADD parent_id integer not null default 0;
-- +migrate Down
ALTER TABLE comments
DROP parent_id;