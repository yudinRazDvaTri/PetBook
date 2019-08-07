
-- +migrate Up
create table if not exists topics (
    id integer not null primary key,
    creator_id integer not null,
    created timestamp not null,
    title text not null,
    description text
);
-- +migrate Down
drop table topics;