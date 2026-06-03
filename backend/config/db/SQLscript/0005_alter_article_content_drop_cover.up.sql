-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0005
-- 描述: 文章内容字段改为 TEXT，删除封面图字段
-- =============================================

-- 将 content 字段类型从 LONGTEXT 改为 TEXT
ALTER TABLE `article` MODIFY `content` TEXT NOT NULL;

-- 删除封面图字段
ALTER TABLE `article` DROP COLUMN `cover`;
