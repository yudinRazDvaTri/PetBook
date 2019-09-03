-- +migrate Up
create table followers
(
    user_id        integer   not null references users (id),
    follower_id    integer   not null references users (id)
);
ALTER TABLE followers ADD CONSTRAINT follow UNIQUE (user_id, follower_id);
-- +migrate Down
drop table followers;