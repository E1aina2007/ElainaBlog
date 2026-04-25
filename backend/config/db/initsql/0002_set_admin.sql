-- =============
-- 初始化管理员账户
-- 密码需要在应用层用 bcrypt 加密后替换下方占位符
-- =============

USE ElainaBlog;

INSERT IGNORE INTO `user` (`username`, `password`, `email`, `is_admin`)
VALUES ('admin', '$2a$10$REPLACE_WITH_BCRYPT_HASH', 'admin@elainaweb.com', 1);
