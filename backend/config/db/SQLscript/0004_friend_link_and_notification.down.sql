-- =============================================
-- ElainaBlog 数据库回滚脚本
-- 版本: 0004
-- 描述: 删除友情链接表和通知表
-- =============================================

DROP TABLE IF EXISTS `notification`;
DROP TABLE IF EXISTS `friend_link`;
