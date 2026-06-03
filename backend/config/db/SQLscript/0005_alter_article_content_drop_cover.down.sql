-- =============================================
-- ElainaBlog 数据库回滚脚本
-- 版本: 0005
-- 描述: 恢复封面图字段，内容字段改回 LONGTEXT
-- =============================================

-- 恢复封面图字段
ALTER TABLE `article` ADD COLUMN `cover` VARCHAR(255) NOT NULL DEFAULT '' AFTER `content`;

-- 将 content 字段类型从 TEXT 改回 LONGTEXT
ALTER TABLE `article` MODIFY `content` LONGTEXT NOT NULL;
