
-- +migrate Up
ALTER TABLE users
    ALTER COLUMN pet_or_vet DROP NOT NULL;
-- +migrate Down