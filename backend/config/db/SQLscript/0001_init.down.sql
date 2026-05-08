-- =============================================
-- ElainaBlog 数据库回滚脚本
-- 版本: 0001
-- 描述: 按外键依赖逆序删除所有核心表
-- =============================================

-- 删除关联表
DROP TABLE IF EXISTS `comment`;

-- 删除主表
DROP TABLE IF EXISTS `article`;
DROP TABLE IF EXISTS `tag`;
DROP TABLE IF EXISTS `category`;
DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `message`;
