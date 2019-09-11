-- +migrate Up
create table if not exists comments
(
    comment_id   serial    not null,
    topic_id     integer   not null references topics (topic_id),
    user_id      integer   not null references users (id),
    created_time TIMESTAMP not null default CURRENT_TIMESTAMP,
    content      text      not null,
    primary key (comment_id)
);
-- +migrate Down
drop table comments;