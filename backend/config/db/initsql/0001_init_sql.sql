-- =============
-- 初始化数据库，设置各表以及字段
-- =============

CREATE DATABASE IF NOT EXISTS ElainaBlog
DEFAULT CHARSET utf8mb4
DEFAULT COLLATE utf8mb4_0900_ai_ci;

USE ElainaBlog;

-- 用户表
CREATE TABLE IF NOT EXISTS `user` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `username`   VARCHAR(50)  NOT NULL,
    `password`   VARCHAR(255) NOT NULL,
    `email`      VARCHAR(100) NOT NULL DEFAULT '',
    `avatar`     VARCHAR(255) NOT NULL DEFAULT '',
    `is_admin`   TINYINT(1)   NOT NULL DEFAULT 0,
    `is_deleted` TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_username` (`username`),
    KEY `idx_email` (`email`),
    KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

-- 分类表
CREATE TABLE IF NOT EXISTS `category` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(50)  NOT NULL,
    `is_deleted` TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`),
    KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='分类表';

-- 文章表
CREATE TABLE IF NOT EXISTS `article` (
    `id`          BIGINT        NOT NULL AUTO_INCREMENT,
    `user_id`     BIGINT        NOT NULL,
    `category_id` BIGINT                 DEFAULT NULL,
    `title`       VARCHAR(200)  NOT NULL,
    `summary`     VARCHAR(500)  NOT NULL DEFAULT '',
    `content`     LONGTEXT      NOT NULL,
    `cover`       VARCHAR(255)  NOT NULL DEFAULT '',
    `is_top`      TINYINT(1)    NOT NULL DEFAULT 0   COMMENT '是否置顶',
    `is_draft`    TINYINT(1)    NOT NULL DEFAULT 0   COMMENT '是否草稿',
    `view_count`  INT           NOT NULL DEFAULT 0,
    `is_deleted`  TINYINT(1)    NOT NULL DEFAULT 0,
    `created_at`  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_is_deleted` (`is_deleted`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_category_id` (`category_id`),
    CONSTRAINT `fk_article_user`     FOREIGN KEY (`user_id`)     REFERENCES `user`(`id`)     ON DELETE CASCADE,
    CONSTRAINT `fk_article_category` FOREIGN KEY (`category_id`) REFERENCES `category`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='文章表';

-- 评论表（仅一级评论）
CREATE TABLE IF NOT EXISTS `comment` (
    `id`         BIGINT        NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT        NOT NULL,
    `user_id`    BIGINT        NOT NULL,
    `content`    TEXT          NOT NULL,
    `is_deleted` TINYINT(1)    NOT NULL DEFAULT 0,
    `created_at` DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_is_deleted` (`is_deleted`),
    KEY `idx_article_id` (`article_id`),
    CONSTRAINT `fk_comment_article` FOREIGN KEY (`article_id`) REFERENCES `article`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comment_user`    FOREIGN KEY (`user_id`)    REFERENCES `user`(`id`)    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评论表';

-- 站点配置表（key-value，存储网站信息、备案号、社交链接等）
CREATE TABLE IF NOT EXISTS `site_config` (
    `key`        VARCHAR(100) NOT NULL,
    `value`      TEXT         NOT NULL,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='站点配置表';

-- 初始化默认站点配置
INSERT IGNORE INTO `site_config` (`key`, `value`) VALUES
    ('title',         '博客标题'),
    ('description',   '博客描述'),
    ('nickname',      '昵称'),
    ('job',           '工作'),
    ('address',       '地址'),
    ('email',         'xxx@qq.com'),
    ('logo_path',     ''),
    ('full_logo_path',''),
    ('icp_number',    ''),
    ('police_number', ''),
    ('bili_url',      ''),
    ('github_url',    '');

