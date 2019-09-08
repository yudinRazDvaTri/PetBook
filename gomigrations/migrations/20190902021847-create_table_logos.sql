
-- +migrate Up
Create table if not exists logos
(
	logo_id serial not null,
	user_id int
		constraint logos_users_id_fk
			references users,
	logo_path varchar,
	created_time TIMESTAMP not null default CURRENT_TIMESTAMP
);

create unique index logos_media_id_uindex
	on logos (logo_id);

create unique index logos_logo_path_uindex
	on logos (logo_path);

alter table logos
	add constraint logos_pk
		primary key (logo_id);


-- +migrate Down
DROP TABLE logos;