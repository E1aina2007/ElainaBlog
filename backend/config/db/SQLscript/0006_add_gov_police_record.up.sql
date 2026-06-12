-- +migrate Up
-- 在 site_config 表中添加公安备案号初始记录
INSERT INTO site_config (key_name, value) VALUES ('gov_police_record', '')
ON DUPLICATE KEY UPDATE value = VALUES(value);
