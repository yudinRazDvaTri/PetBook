-- +migrate Up
create table if not exists refresh_tokens
(
    id             serial    not null,
    user_id        integer   not null references users (id),
    token_string   text      not null unique,
    last_update_at timestamp not null default CURRENT_TIMESTAMP,
    primary key (id)
);

-- +migrate Down
drop table refresh_tokens;
