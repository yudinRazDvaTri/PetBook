-- +migrate Up
ALTER TABLE pets
    ADD constraint pets_fk
    FOREIGN KEY (user_id) references users(id);
ALTER TABLE vets
    ADD constraint vets_fk
    FOREIGN KEY (user_id) references users(id);
-- +migrate Down
ALTER TABLE pets DROP CONSTRAINT pets_fk;
ALTER TABLE vets DROP CONSTRAINT vets_fk;
