
-- +migrate Up
DROP TABLE messages;
CREATE TABLE IF NOT EXISTS messages(
id serial primary key,
to_id integer not null references users (id),
from_id integer not null references users (id),
text text not null ,
created_at TIMESTAMP not null
);
-- +migrate Down
DROP TABLE messages;
CREATE TABLE IF NOT EXISTS messages(
id integer not null primary key,
to_id integer not null references users (id),
from_id integer not null references users (id),
text text not null ,
created_at date not null
);

