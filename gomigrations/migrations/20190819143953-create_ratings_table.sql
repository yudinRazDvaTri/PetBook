
-- +migrate Up
create table if not exists ratings
(
    comment_id   integer   not null references comments (comment_id),
    user_id      integer   not null references users (id)
);
ALTER TABLE ratings ADD CONSTRAINT rating UNIQUE (comment_id, user_id);
-- +migrate Down
drop table ratings;