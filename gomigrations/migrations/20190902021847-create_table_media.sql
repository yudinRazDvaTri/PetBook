
-- +migrate Up
Create table if not exists media
(
	media_id serial not null,
	user_id int
		constraint media_users_id_fk
			references users,
	logo_path varchar,
	gallery_path varchar,
	docs_path varchar,
	created_time TIMESTAMP not null default CURRENT_TIMESTAMP
);

create unique index media_media_id_uindex
	on media (media_id);

create unique index media_logo_path_uindex
	on media (logo_path);

alter table media
	add constraint media_pk
		primary key (media_id);


-- +migrate Down
DROP TABLE media;