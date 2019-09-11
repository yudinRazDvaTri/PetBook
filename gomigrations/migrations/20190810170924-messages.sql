-- +migrate Up
CREATE TABLE IF NOT EXISTS messages(
id integer not null primary key,
to_id integer not null,
from_id integer not null,
text text not null ,
created_at date not null ,
read boolean not null 
);
-- +migrate Down
DROP TABLE messages;
