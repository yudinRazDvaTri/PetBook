
-- +migrate Up
ALTER TABLE comments
    ADD constraint topics_fkey
    FOREIGN KEY (topic_id) references topics(id);
-- +migrate Down
ALTER TABLE comments DROP CONSTRAINT topics_fkey;