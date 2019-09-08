
-- +migrate Up
alter table users alter column pet_or_vet set not null;

-- +migrate Down
 alter column pet_or_vet DROP NOT NULL;