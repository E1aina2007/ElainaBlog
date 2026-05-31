-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0004
-- 描述: 创建友情链接表和通知表
-- =============================================

-- 友情链接表：存储博主的友情链接
CREATE TABLE IF NOT EXISTS `friend_link` (
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(100) NOT NULL COMMENT '站点名称',
    `url`         VARCHAR(500) NOT NULL COMMENT '站点链接',
    `avatar`      VARCHAR(500) NOT NULL DEFAULT '' COMMENT '头像/Logo URL',
    `description` VARCHAR(300) NOT NULL DEFAULT '' COMMENT '站点描述',
    `sort_order`  INT          NOT NULL DEFAULT 0 COMMENT '排序权重，越大越靠前',
    `is_deleted`  TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_is_deleted` (`is_deleted`),
    KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 通知表：存储用户通知（评论提醒、留言提醒等）
CREATE TABLE IF NOT EXISTS `notification` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `user_id`    BIGINT       NOT NULL COMMENT '接收者ID',
    `type`       VARCHAR(20)  NOT NULL COMMENT '类型: comment/message',
    `title`      VARCHAR(200) NOT NULL COMMENT '通知标题',
    `content`    VARCHAR(500) NOT NULL DEFAULT '' COMMENT '通知内容摘要',
    `target_id`  BIGINT       NOT NULL DEFAULT 0 COMMENT '关联目标ID（文章ID/留言ID）',
    `is_read`    TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否已读',
    `is_deleted` TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_is_read` (`is_read`),
    KEY `idx_is_deleted` (`is_deleted`),
    CONSTRAINT `fk_notification_user` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
