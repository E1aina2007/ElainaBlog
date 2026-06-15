-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0009 (回滚)
-- 描述: 移除评论回复目标评论ID和内容字段
-- =============================================

ALTER TABLE comment DROP COLUMN reply_to_content;
ALTER TABLE comment DROP COLUMN reply_to_comment_id;
