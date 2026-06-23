-- 回滚: 删除文章关键词标签字段
ALTER TABLE article DROP COLUMN tags;
