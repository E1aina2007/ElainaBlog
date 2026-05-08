-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0002
-- 描述: 创建留言表（作者页面留言板功能）
-- =============================================

-- 留言表：用户在作者页面的留言
CREATE TABLE IF NOT EXISTS `message` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `user_id`    BIGINT       NOT NULL,                          -- 留言者ID
    `content`    TEXT         NOT NULL,                          -- 留言内容
    `is_deleted` TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_is_deleted` (`is_deleted`),
    CONSTRAINT `fk_message_user` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
