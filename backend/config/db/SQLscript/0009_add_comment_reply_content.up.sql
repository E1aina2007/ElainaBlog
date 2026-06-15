-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0009
-- 描述: 添加评论回复目标评论ID和内容字段
-- =============================================

-- 添加回复目标评论ID字段
ALTER TABLE comment ADD COLUMN reply_to_comment_id BIGINT DEFAULT NULL AFTER reply_to_username;
-- 添加回复目标评论内容字段（反序列化快照）
ALTER TABLE comment ADD COLUMN reply_to_content TEXT DEFAULT NULL AFTER reply_to_comment_id;

-- 添加回复目标评论索引
CREATE INDEX idx_comment_reply_to_comment ON comment(reply_to_comment_id);
