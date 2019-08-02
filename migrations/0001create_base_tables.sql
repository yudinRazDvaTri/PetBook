-- +migrate Up
CREATE type role AS ENUM('pet','vet');
CREATE type sex AS ENUM('male','female');
CREATE type kind_of_animal AS ENUM('dog','cat','fish','bird','rodent');
CREATE TABLE IF NOT EXISTS users(
id serial not null PRIMARY KEY,
email text not null unique,
first_name text not null,
last_name text not null ,
login text unique not null ,
pet_or_vet role not null ,
password text not null
);

CREATE table if not exists pets(
user_id integer not null primary key,
name text not null,
ade integer not null,
animal_type kind_of_animal not null,
breed text not null,
description text default 'There is no description yet',
height float not null,
weight float not null,
gender sex not null
);

CREATE table if not exists vets (
user_id integer not null primary key,
experience float not null,
description text default 'There is no description yet'
);
-- +migrate Down
DROP TABLE users;
DROP TABLE pets;
DROP TABLE vets;
DROP type role;
DROP type sex;
DROP type kind_of_animal;