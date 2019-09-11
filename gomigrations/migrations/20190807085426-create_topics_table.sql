-- +migrate Up
create table if not exists topics
(
    topic_id     serial    not null,
    user_id      integer   not null references users (id),
    created_time TIMESTAMP not null default CURRENT_TIMESTAMP,
    title        text      not null,
    description  text,
    primary key (topic_id)
);
-- +migrate Down
drop table topics;