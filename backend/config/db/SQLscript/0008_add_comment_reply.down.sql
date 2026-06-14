-- =============================================
-- ElainaBlog 数据库回滚脚本
-- 版本: 0008
-- 描述: 删除评论回复功能字段
-- =============================================

-- 删除回复目标用户索引和字段
ALTER TABLE comment DROP INDEX idx_comment_reply_to;
ALTER TABLE comment DROP COLUMN reply_to_user_id;
ALTER TABLE comment DROP COLUMN reply_to_username;
