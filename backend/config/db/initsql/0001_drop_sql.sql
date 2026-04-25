-- =============
-- 回滚脚本：按外键依赖的逆序删除所有表
-- =============

USE ElainaBlog;

DROP TABLE IF EXISTS `comment`;
DROP TABLE IF EXISTS `article_tag`;
DROP TABLE IF EXISTS `article`;
DROP TABLE IF EXISTS `tag`;
DROP TABLE IF EXISTS `category`;
DROP TABLE IF EXISTS `user`;
