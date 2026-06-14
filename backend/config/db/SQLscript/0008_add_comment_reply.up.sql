-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0008
-- 描述: 添加评论回复功能字段
-- =============================================

-- 添加回复目标用户ID和用户名字段
ALTER TABLE comment ADD COLUMN reply_to_user_id BIGINT DEFAULT NULL AFTER user_id;
ALTER TABLE comment ADD COLUMN reply_to_username VARCHAR(50) DEFAULT NULL AFTER reply_to_user_id;

-- 添加回复目标用户索引
CREATE INDEX idx_comment_reply_to ON comment(reply_to_user_id);
