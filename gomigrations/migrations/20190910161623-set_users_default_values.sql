-- +migrate Up
ALTER table users
    ALTER column firstname SET DEFAULT '',
    ALTER column lastname SET DEFAULT '';

-- +migrate Down
