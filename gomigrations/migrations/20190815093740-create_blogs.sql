
-- +migrate Up
create table if not exists blog (
    blog_id serial not null,
    user_id integer not null references users(id),
    created_time TIMESTAMP not null default CURRENT_TIMESTAMP,
    content text not null,
    primary key(blog_id));

-- +migrate Down
drop table comments;