ALTER TABLE comment DROP INDEX idx_comment_reply_to;
ALTER TABLE comment DROP COLUMN reply_to_user_id;
ALTER TABLE comment DROP COLUMN reply_to_username;
