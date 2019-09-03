-- +migrate Up
ALTER TABLE vets
    ADD constraint users_fk
        FOREIGN KEY (user_id) references users (id);
-- +migrate Down
ALTER TABLE vets
    DROP CONSTRAINT users_fk;