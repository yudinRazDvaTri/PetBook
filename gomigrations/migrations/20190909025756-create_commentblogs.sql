
-- +migrate Up
create table if not exists commentblog (
    commentblog_id serial not null,
    blog_id integer not null references blog (blog_id),
    user_id integer not null references users (id),
    created_time TIMESTAMP not null default CURRENT_TIMESTAMP,
    content text not null,
    primary key(commentblog_id));

-- +migrate Down
drop table comments;