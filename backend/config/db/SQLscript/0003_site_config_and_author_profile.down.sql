-- =============================================
-- ElainaBlog 数据库回滚脚本
-- 版本: 0003
-- 描述: 删除站点配置表和作者信息表
-- =============================================

DROP TABLE IF EXISTS `author_profile`;
DROP TABLE IF EXISTS `site_config`;
