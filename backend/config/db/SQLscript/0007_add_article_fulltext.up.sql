-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0007
-- 描述: 添加文章全文搜索索引
-- =============================================

-- 为文章标题和摘要添加全文索引（使用 ngram 解析器支持中文）
ALTER TABLE article ADD FULLTEXT INDEX ft_title_summary (title, summary) WITH PARSER ngram;
