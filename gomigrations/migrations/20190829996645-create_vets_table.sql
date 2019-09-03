-- +migrate Up
CREATE type rank AS ENUM ('A','B','C','D');
CREATE type class AS ENUM ('zoologist','handler','paramedic','surgeon');
CREATE table if not exists vets
(
    user_id     integer        not null primary key,
    name        text           not null,
    qualification         class        not null,
    surname text not null,
    category       rank           not null,
    certificates text
);
-- +migrate Down
DROP TABLE vets;
DROP type rank;
DROP type class;