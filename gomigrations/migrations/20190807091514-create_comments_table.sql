
-- +migrate Up
create table if not exists comments (
    id integer not null primary key,
    topic_id integer not null,
    creator_id integer not null,
    created timestamp not null,
    edited timestamp,
    content text not null
);
-- +migrate Down
drop table comments;