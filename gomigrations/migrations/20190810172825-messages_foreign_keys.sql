
-- +migrate Up
ALTER TABLE messages
    ADD constraint user_to_fk
    FOREIGN KEY (to_id) references users(id);
ALTER TABLE messages
    ADD constraint user_from_fk
    FOREIGN KEY (from_id) references users(id);
-- +migrate Down
ALTER TABLE messages
 DROP CONSTRAINT user_to_fk;
 DROP CONSTRAINT user_from_fk;