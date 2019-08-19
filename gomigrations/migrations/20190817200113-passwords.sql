-- +migrate Up
CREATE TABLE IF NOT EXISTS passwords
(
    id              serial not null,
    user_id         int    not null references users(id),
    password_string text   not null
);
-- +migrate Down
DROP TABLE passwords;
