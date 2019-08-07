
-- +migrate Up
ALTER TABLE topics
    ADD constraint users_fkey
    FOREIGN KEY (creator_id) references users(id);
-- +migrate Down
ALTER TABLE topics DROP CONSTRAINT users_fkey;