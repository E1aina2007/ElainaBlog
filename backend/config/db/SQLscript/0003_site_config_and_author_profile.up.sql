-- =============================================
-- ElainaBlog 数据库迁移脚本
-- 版本: 0003
-- 描述: 创建站点配置表和作者信息表
-- =============================================

-- 站点配置表：存储站点名称、标题、问候语、随机语句等可配置项
CREATE TABLE IF NOT EXISTS `site_config` (
    `id`         BIGINT       NOT NULL AUTO_INCREMENT,
    `key_name`   VARCHAR(100) NOT NULL COMMENT '配置键',
    `value`      TEXT         NOT NULL COMMENT '配置值',
    `is_deleted` TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_key_name` (`key_name`),
    KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 作者信息表：存储作者个人信息（单行设计，仅维护一条记录）
CREATE TABLE IF NOT EXISTS `author_profile` (
    `id`                  BIGINT       NOT NULL AUTO_INCREMENT,
    `nickname`            VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar`              VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像URL',
    `background`          VARCHAR(255) NOT NULL DEFAULT '' COMMENT '背景图URL',
    `signature`           VARCHAR(255) NOT NULL DEFAULT '' COMMENT '个性签名',
    `location`            VARCHAR(100) NOT NULL DEFAULT '' COMMENT '所在城市',
    `occupation`          VARCHAR(100) NOT NULL DEFAULT '' COMMENT '职业',
    `school`              VARCHAR(100) NOT NULL DEFAULT '' COMMENT '院校',
    `major`               VARCHAR(100) NOT NULL DEFAULT '' COMMENT '专业',
    `email`               VARCHAR(100) NOT NULL DEFAULT '' COMMENT '邮箱',
    `wechat`              VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '微信',
    `bio`                 TEXT         NOT NULL COMMENT '个人简介',
    `tech_stack_frontend` TEXT         NOT NULL COMMENT '前端技术栈JSON数组',
    `tech_stack_backend`  TEXT         NOT NULL COMMENT '后端技术栈JSON数组',
    `tech_stack_engineering` TEXT      NOT NULL COMMENT '工程化技术栈JSON数组',
    `social_github`       VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'GitHub链接',
    `social_bilibili`     VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Bilibili链接',
    `is_deleted`          TINYINT(1)   NOT NULL DEFAULT 0,
    `created_at`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 插入默认站点配置
INSERT INTO `site_config` (`key_name`, `value`) VALUES
    ('site_name', 'My Blog'),
    ('site_title', 'ElainaBlog'),
    ('greeting', 'Hello!'),
    ('hero_title', 'Welcome to my blog'),
    ('icp_beian', ''),
    ('quotes', '["生活不止眼前的苟且","诗和远方就在脚下","代码改变世界","保持热爱，奔赴山海"]');

-- 插入默认作者信息
INSERT INTO `author_profile` (`nickname`, `avatar`, `background`, `signature`, `location`, `occupation`, `school`, `major`, `email`, `wechat`, `bio`, `tech_stack_frontend`, `tech_stack_backend`, `tech_stack_engineering`, `social_github`, `social_bilibili`) VALUES
    ('Author', '/author/avatar.jpg', '/author/background.jpg', 'Hello World', '', '', '', '', '', '', 'This is a blog powered by ElainaBlog.', '["HTML5","CSS 3","JavaScript/TypeScript","Vue 3"]', '["Go","MySQL","Redis"]', '["Docker","Git","Nginx","Vite"]', '', '');
