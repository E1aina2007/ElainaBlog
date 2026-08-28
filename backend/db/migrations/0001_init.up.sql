-- =============================================
-- ElainaBlog 数据库初始化脚本
-- 由 golang-migrate 容器执行，后端不再内置迁移
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

-- 将 content 字段类型从 LONGTEXT 改为 TEXT
ALTER TABLE `article` MODIFY `content` TEXT NOT NULL;

-- 删除封面图字段
ALTER TABLE `article` DROP COLUMN `cover`;

-- 在 site_config 表中添加公安备案号初始记录
INSERT INTO site_config (key_name, value) VALUES ('gov_police_record', '')
ON DUPLICATE KEY UPDATE value = VALUES(value);

-- 为文章标题和摘要添加全文索引（使用 ngram 解析器支持中文）
ALTER TABLE article ADD FULLTEXT INDEX ft_title_summary (title, summary) WITH PARSER ngram;

-- 添加回复目标用户ID和用户名字段
ALTER TABLE comment ADD COLUMN reply_to_user_id BIGINT DEFAULT NULL AFTER user_id;
ALTER TABLE comment ADD COLUMN reply_to_username VARCHAR(50) DEFAULT NULL AFTER reply_to_user_id;

-- 添加回复目标用户索引
CREATE INDEX idx_comment_reply_to ON comment(reply_to_user_id);

-- 添加回复目标评论ID字段
ALTER TABLE comment ADD COLUMN reply_to_comment_id BIGINT DEFAULT NULL AFTER reply_to_username;
-- 添加回复目标评论内容字段（反序列化快照）
ALTER TABLE comment ADD COLUMN reply_to_content TEXT DEFAULT NULL AFTER reply_to_comment_id;

-- 添加回复目标评论索引
CREATE INDEX idx_comment_reply_to_comment ON comment(reply_to_comment_id);

-- 添加 tags 字段，存储逗号分隔的关键词
ALTER TABLE article ADD COLUMN tags VARCHAR(500) DEFAULT '' AFTER is_draft;

-- 分类表添加置顶标记
ALTER TABLE category ADD COLUMN is_top TINYINT(1) NOT NULL DEFAULT 0 AFTER name;
