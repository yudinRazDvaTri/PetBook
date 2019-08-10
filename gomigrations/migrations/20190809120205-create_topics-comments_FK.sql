
-- +migrate Up
ALTER TABLE topics
    ADD constraint topics_asnwer_id_fkey
    FOREIGN KEY (answer_id) references comments(comment_id);
-- +migrate Down
ALTER TABLE topics DROP CONSTRAINT topics_asnwer_id_fkey;