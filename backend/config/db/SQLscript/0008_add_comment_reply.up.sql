ALTER TABLE comment ADD COLUMN reply_to_user_id BIGINT DEFAULT NULL AFTER user_id;
ALTER TABLE comment ADD COLUMN reply_to_username VARCHAR(50) DEFAULT NULL AFTER reply_to_user_id;
CREATE INDEX idx_comment_reply_to ON comment(reply_to_user_id);
