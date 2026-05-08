-- =============================================
-- ElainaBlog 数据库初始化脚本
-- 版本: 0001
-- 描述: 创建博客系统核心表结构
-- =============================================

-- 用户表：存储用户信息
CREATE TABLE IF NOT EXISTS `user` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `username`   VARCHAR(50)  NOT NULL,                          -- 用户名
    `password`   VARCHAR(255) NOT NULL,                          -- 密码（bcrypt加密）
    `email`      VARCHAR(100) NOT NULL DEFAULT '',               -- 邮箱
    `avatar`     VARCHAR(255) NOT NULL DEFAULT '',               -- 头像URL
    `is_admin`   TINYINT(1)   NOT NULL DEFAULT 0,               -- 是否管理员
    `is_deleted` TINYINT(1)   NOT NULL DEFAULT 0,               -- 软删除标记
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_username` (`username`),
    KEY `idx_email` (`email`),
    KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 分类表：文章分类
CREATE TABLE IF NOT EXISTS `category` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(50)  NOT NULL,                          -- 分类名称
    `is_deleted` TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`),
    KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 文章表：存储博客文章
CREATE TABLE IF NOT EXISTS `article` (
    `id`          BIGINT        NOT NULL AUTO_INCREMENT,
    `user_id`     BIGINT        NOT NULL,                       -- 作者ID
    `category_id` BIGINT                 DEFAULT NULL,          -- 分类ID
    `title`       VARCHAR(200)  NOT NULL,                       -- 标题
    `summary`     VARCHAR(500)  NOT NULL DEFAULT '',             -- 摘要
    `content`     LONGTEXT      NOT NULL,                       -- 内容（Markdown）
    `cover`       VARCHAR(255)  NOT NULL DEFAULT '',             -- 封面图URL
    `is_top`      TINYINT(1)    NOT NULL DEFAULT 0,             -- 是否置顶
    `is_draft`    TINYINT(1)    NOT NULL DEFAULT 0,             -- 是否草稿
    `view_count`  INT           NOT NULL DEFAULT 0,             -- 阅读量
    `is_deleted`  TINYINT(1)    NOT NULL DEFAULT 0,
    `created_at`  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_is_deleted` (`is_deleted`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_category_id` (`category_id`),
    CONSTRAINT `fk_article_user`     FOREIGN KEY (`user_id`)     REFERENCES `user`(`id`)     ON DELETE CASCADE,
    CONSTRAINT `fk_article_category` FOREIGN KEY (`category_id`) REFERENCES `category`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 评论表：文章评论（仅支持一级评论）
CREATE TABLE IF NOT EXISTS `comment` (
    `id`         BIGINT        NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT        NOT NULL,                        -- 文章ID
    `user_id`    BIGINT        NOT NULL,                        -- 评论者ID
    `content`    TEXT          NOT NULL,                        -- 评论内容
    `is_deleted` TINYINT(1)    NOT NULL DEFAULT 0,
    `created_at` DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_is_deleted` (`is_deleted`),
    KEY `idx_article_id` (`article_id`),
    CONSTRAINT `fk_comment_article` FOREIGN KEY (`article_id`) REFERENCES `article`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comment_user`    FOREIGN KEY (`user_id`)    REFERENCES `user`(`id`)    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 标签表：文章标签
CREATE TABLE IF NOT EXISTS `tag` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(50)  NOT NULL,                          -- 标签名称
    `is_deleted` TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`),
    KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;