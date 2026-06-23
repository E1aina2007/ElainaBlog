-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0010
-- 描述: 添加文章关键词标签字段
-- =============================================

-- 添加 tags 字段，存储逗号分隔的关键词
ALTER TABLE article ADD COLUMN tags VARCHAR(500) DEFAULT '' AFTER is_draft;
