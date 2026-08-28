-- =============================================
-- ElainaBlog 数据库回滚脚本（合并版）
-- 按变更与依赖的逆序回滚全部结构
-- =============================================

ALTER TABLE category DROP COLUMN is_top;

ALTER TABLE article DROP COLUMN tags;

ALTER TABLE comment DROP COLUMN reply_to_content;
ALTER TABLE comment DROP COLUMN reply_to_comment_id;

ALTER TABLE comment DROP INDEX idx_comment_reply_to;
ALTER TABLE comment DROP COLUMN reply_to_user_id;
ALTER TABLE comment DROP COLUMN reply_to_username;

ALTER TABLE article DROP INDEX ft_title_summary;

DELETE FROM site_config WHERE key_name = 'gov_police_record';

-- 恢复封面图字段
ALTER TABLE `article` ADD COLUMN `cover` VARCHAR(255) NOT NULL DEFAULT '' AFTER `content`;

-- 将 content 字段类型从 TEXT 改回 LONGTEXT
ALTER TABLE `article` MODIFY `content` LONGTEXT NOT NULL;

DROP TABLE IF EXISTS `notification`;
DROP TABLE IF EXISTS `friend_link`;

DROP TABLE IF EXISTS `author_profile`;
DROP TABLE IF EXISTS `site_config`;

DROP TABLE IF EXISTS `message`;

DROP TABLE IF EXISTS `comment`;
DROP TABLE IF EXISTS `article`;
DROP TABLE IF EXISTS `tag`;
DROP TABLE IF EXISTS `category`;
DROP TABLE IF EXISTS `user`;
