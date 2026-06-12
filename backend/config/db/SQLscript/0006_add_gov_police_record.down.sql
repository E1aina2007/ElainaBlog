-- +migrate Down
-- 删除公安备案号配置记录
DELETE FROM site_config WHERE key_name = 'gov_police_record';
