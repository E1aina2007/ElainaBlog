# 图片资源管理方案

> 解决文章删除后图片残留、防止恶意上传占用服务器资源

---

## 问题分析

### 1. 孤儿图片问题

**现状**：
- 图片 URL 以 Markdown 语法嵌入文章 `content` 字段（`![](/uploads/2026-06-14/xxx.jpg)`）
- 文章删除使用软删除（`is_deleted=1`），content 仍保留在数据库
- 系统没有图片引用追踪机制
- 删除文章时只级联删除评论，不处理图片文件
- 上传后未保存到文章的图片也会残留（如编辑草稿时上传但未引用）

**后果**：
- 服务器磁盘空间持续增长，无法回收
- 孤儿图片无法被用户访问，造成资源浪费

### 2. 恶意上传风险

**现状防护**：
- ✅ 需要 JWT 登录才能上传
- ✅ 文件大小限制（20MB）
- ✅ 扩展名白名单 + MIME 检测

**缺失防护**：
- ❌ 无上传频率限制（Rate Limiting）
- ❌ 无用户存储配额
- ❌ 无全局存储上限监控
- ❌ 草稿中上传但未引用的图片无法追踪

---

## 解决方案

### 方案一：孤儿图片定时清理（推荐）

**原理**：定期扫描 uploads 目录，对比所有有效文章中的图片引用，删除未被引用的文件。

**优点**：
- 不侵入业务逻辑，无需修改删除文章的代码
- 可清理各种来源的孤儿图片（草稿删除、编辑时替换等）
- 实现简单，一个定时任务即可

**实现步骤**：

#### 1. 新增图片清理服务

```go
// backend/internal/upload/cleanup.go

package upload

import (
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "log"
)

// ImageCleanup 图片清理服务
type ImageCleanup struct {
    uploadPath string
    repo       ArticleRepository // 需要查询所有文章 content
}

// ArticleRepository 文章仓库接口（仅需要的方法）
type ArticleRepository interface {
    GetAllActiveContents() ([]string, error) // 获取所有未删除文章的 content
}

// CleanupOrphanImages 清理孤儿图片
func (s *ImageCleanup) CleanupOrphanImages() (int, error) {
    // 1. 获取所有有效文章的 content
    contents, err := s.repo.GetAllActiveContents()
    if err != nil {
        return 0, err
    }

    // 2. 提取所有被引用的图片 URL
    referencedImages := make(map[string]bool)
    imgRegex := regexp.MustCompile(`!\[.*?\]\((/uploads/\d{4}-\d{2}-\d{2}/[^)]+)\)`)
    for _, content := range contents {
        matches := imgRegex.FindAllStringSubmatch(content, -1)
        for _, match := range matches {
            if len(match) > 1 {
                // 转换为文件系统路径
                relPath := strings.TrimPrefix(match[1], "/uploads/")
                referencedImages[relPath] = true
            }
        }
    }

    // 3. 扫描 uploads 目录
    deletedCount := 0
    err = filepath.Walk(s.uploadPath, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }

        // 跳过 avatars 目录
        if strings.Contains(path, "avatars") {
            return nil
        }

        // 获取相对路径
        relPath, _ := filepath.Rel(s.uploadPath, path)
        relPath = filepath.ToSlash(relPath)

        // 4. 不在引用集合中则删除
        if !referencedImages[relPath] {
            if err := os.Remove(path); err != nil {
                log.Printf("[ImageCleanup] 删除失败: %s, error: %v", path, err)
            } else {
                deletedCount++
                log.Printf("[ImageCleanup] 已删除孤儿图片: %s", relPath)
            }
        }
        return nil
    })

    return deletedCount, err
}
```

#### 2. 新增 Repository 方法

```go
// backend/internal/article/repository.go

// GetAllActiveContents 获取所有未删除文章的 content
func (r *articleRepository) GetAllActiveContents() ([]string, error) {
    var contents []string
    err := r.db.Raw("SELECT content FROM article WHERE is_deleted = 0").Scan(&contents).Error
    return contents, err
}
```

#### 3. 注册定时任务

```go
// backend/cmd/handler.go

// 在 InitScheduledTasks 中添加
func InitScheduledTasks() {
    // ... 现有任务 ...

    // 每天凌晨 3 点清理孤儿图片
    scheduler.Every(1).Day().At("03:00").Do(func() {
        count, err := imageCleanup.CleanupOrphanImages()
        if err != nil {
            log.Printf("[Scheduler] 图片清理失败: %v", err)
        } else {
            log.Printf("[Scheduler] 图片清理完成，删除 %d 个孤儿图片", count)
        }
    })
}
```

---

### 方案二：上传频率限制 + 用户配额

**原理**：限制每个用户的上传频率和总存储空间，防止恶意上传。

#### 1. 上传频率限制（Redis 计数器）

```go
// backend/internal/upload/middleware.go

// UploadRateLimit 上传频率限制中间件
func UploadRateLimit(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("userID")
        key := fmt.Sprintf("upload:rate:%s:%d", userID, time.Now().Unix()/int64(window.Seconds()))

        count, err := rdb.Incr(c, key).Result()
        if err != nil {
            c.Next()
            return
        }

        if count == 1 {
            rdb.Expire(c, key, window)
        }

        if count > int64(limit) {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": fmt.Sprintf("上传过于频繁，每 %v 最多上传 %d 张", window, limit),
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

**配置示例**：
```yaml
# config.yaml
upload:
  rate_limit: 20        # 每小时最多上传 20 张
  rate_window: 3600     # 时间窗口（秒）
```

#### 2. 用户存储配额

```go
// backend/internal/upload/quota.go

// CheckUserQuota 检查用户存储配额
func CheckUserQuota(db *gorm.DB, userID uint, newFileSize int64, quota int64) error {
    var totalSize int64
    err := db.Raw(`
        SELECT COALESCE(SUM(file_size), 0)
        FROM user_uploads
        WHERE user_id = ?
    `, userID).Scan(&totalSize).Error
    if err != nil {
        return err
    }

    if totalSize+newFileSize > quota {
        return fmt.Errorf("存储空间不足，已使用 %dMB，配额 %dMB", totalSize/1024/1024, quota/1024/1024)
    }
    return nil
}
```

**新增数据库表**：
```sql
-- 记录用户上传的文件信息
CREATE TABLE user_uploads (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id)
);
```

---

### 方案三：全局存储监控告警

**原理**：监控 uploads 目录总大小，超过阈值时发送告警。

```go
// backend/internal/upload/monitor.go

// StorageMonitor 存储监控
type StorageMonitor struct {
    uploadPath string
    threshold  int64  // 告警阈值（字节）
    notifier   Notifier
}

type Notifier interface {
    Send(message string) error
}

// CheckStorage 检查存储使用情况
func (m *StorageMonitor) CheckStorage() (int64, error) {
    var totalSize int64
    err := filepath.Walk(m.uploadPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            totalSize += info.Size()
        }
        return nil
    })
    if err != nil {
        return 0, err
    }

    if totalSize > m.threshold {
        msg := fmt.Sprintf("⚠️ 存储空间告警\n当前使用: %.2fGB\n阈值: %.2fGB\n路径: %s",
            float64(totalSize)/1024/1024/1024,
            float64(m.threshold)/1024/1024/1024,
            m.uploadPath,
        )
        m.notifier.Send(msg)
    }

    return totalSize, nil
}
```

---

## 实施计划

### 阶段一：基础防护（1-2 天）

| 任务 | 说明 |
|------|------|
| 实现上传频率限制 | Redis 计数器，每小时 20 张 |
| 添加配置项 | `upload.rate_limit`、`upload.rate_window` |
| 注册中间件 | 在上传路由添加限流中间件 |

### 阶段二：孤儿清理（2-3 天）

| 任务 | 说明 |
|------|------|
| 实现 `ImageCleanup` 服务 | 正则提取 + 文件扫描 + 删除 |
| 新增 Repository 方法 | `GetAllActiveContents` |
| 注册定时任务 | 每天凌晨 3 点执行 |
| 添加日志记录 | 记录删除的文件列表 |

### 阶段三：高级防护（3-5 天，可选）

| 任务 | 说明 |
|------|------|
| 用户上传记录表 | `user_uploads` 迁移脚本 |
| 用户配额检查 | 上传前校验 |
| 存储监控告警 | 钉钉/邮件通知 |

---

## 配置参考

```yaml
# config.yaml
upload:
  size: 20                    # 单文件最大 20MB
  path: uploads               # 存储路径
  avatar_path: uploads/avatars
  avatar_size: 5

  # 新增配置
  rate_limit: 20              # 每小时最多上传 20 张
  rate_window: 3600           # 频率限制时间窗口（秒）
  user_quota: 104857600       # 用户存储配额 100MB（可选）
  storage_threshold: 10737418240  # 全局告警阈值 10GB（可选）
  cleanup_cron: "0 3 * * *"   # 孤儿清理 cron 表达式
```

---

## 注意事项

1. **头像目录排除**：清理任务需跳过 `avatars` 目录，头像由用户管理覆盖
2. **软删除文章**：当前文章是软删除，清理任务只扫描 `is_deleted=0` 的文章
3. **并发安全**：清理任务执行期间不应有上传操作，建议凌晨低峰期运行
4. **日志审计**：记录所有删除操作，便于追溯和恢复
5. **备份策略**：清理前可先将孤儿图片移动到临时目录，观察一段时间后再彻底删除

---

## 相关文件

- `backend/internal/upload/controller.go` — 上传控制器
- `backend/internal/upload/storage.go` — 存储逻辑
- `backend/internal/article/repository.go` — 文章 Repository
- `backend/cmd/handler.go` — 定时任务注册
- `backend/config/config.go` — 配置结构体
