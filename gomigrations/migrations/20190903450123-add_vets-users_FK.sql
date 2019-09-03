-- +migrate Up
ALTER TABLE pets
    ADD constraint users_fk
        FOREIGN KEY (user_id) references users (id);
-- +migrate Down
ALTER TABLE pets
    DROP CONSTRAINT users_fk;