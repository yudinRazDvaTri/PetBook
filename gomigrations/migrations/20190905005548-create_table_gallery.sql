
-- +migrate Up
Create table if not exists gallery
(
	file_id serial not null,
	user_id int
		constraint gallery_users_id_fk
			references users,
	file_path varchar,
	created_time TIMESTAMP not null default CURRENT_TIMESTAMP
);

create unique index gallery_file_id_uindex
	on gallery (file_id);

create unique index gallery_file_path_uindex
	on gallery (file_path);

alter table gallery
	add constraint gallery_pk
		primary key (file_id);


-- +migrate Down
DROP TABLE gallery;