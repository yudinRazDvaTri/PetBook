
-- +migrate Up
ALTER TABLE comments
    ADD constraint users_fkey
    FOREIGN KEY (creator_id) references users(id);
-- +migrate Down
ALTER TABLE comments DROP CONSTRAINT users_fkey;