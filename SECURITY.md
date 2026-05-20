# ElainaBlog 安全措施

## 1. 网络与接入层防护

### Cloudflare CDN 防护

项目通过 Cloudflare 免费套餐实现以下防护：

| 措施 | 说明 |
|------|------|
| 源站 IP 隐藏 | 所有流量经 Cloudflare 代理，真实服务器 IP 不暴露 |
| DDoS 基础防护 | 自动检测并缓解常见 DDoS 攻击 |
| 恶意流量拦截 | 基于威胁评分自动拦截可疑请求 |
| HTTPS 强制 | 自动将 HTTP 请求重定向到 HTTPS |

### 防火墙规则

Cloudflare 层配置了以下 WAF 规则：

- 拦截包含 SQL 注入特征的请求（`select%20`、`union%20select`）
- 拦截包含 XSS 特征的请求（`<script`）
- 阻止访问敏感文件（`.env`）
- 对高威胁评分 IP 发起质询

---

## 2. 应用层安全

### SQL 注入防护

后端使用 Go 标准库 `database/sql` 的参数化查询，所有 SQL 语句均使用占位符，不拼接用户输入：

```go
db.QueryRow("SELECT id, username, password FROM user WHERE id = ?", id)
```

### XSS 防护

- Nginx 配置了 `X-Content-Type-Options: nosniff` 防止 MIME 类型嗅探
- Nginx 配置了 `X-XSS-Protection: 1; mode=block` 启用浏览器 XSS 过滤
- 前端使用 Vue 3 模板语法，默认对输出进行转义

### 文件上传安全

- 仅允许特定后缀的文件上传（白名单机制）
- 上传文件使用随机文件名存储
- 上传目录与 Web 可执行目录分离
- 限制上传文件大小为 20MB

### 认证与授权

- 密码使用 bcrypt 算法加盐哈希存储
- JWT Token 用于身份认证，支持 Access Token + Refresh Token 双 Token 机制
- 管理员接口通过中间件进行权限校验

---

## 3. 服务器安全配置

### Nginx 安全头

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
```

### 敏感文件保护

Nginx 配置禁止访问隐藏文件和备份文件：

```nginx
location ~ /\.(?!well-known) {
    deny all;
}
```

### Docker 容器隔离

- MySQL 和 Redis 端口仅绑定 `127.0.0.1`，不暴露到公网
- 各服务运行在独立容器中，通过 Docker 网络通信
- 后端容器使用非 root 用户运行

### 系统防火墙

服务器使用 UFW 防火墙，仅开放必要端口：

| 端口 | 协议 | 用途 |
|------|------|------|
| 22 | TCP | SSH 远程管理 |
| 80 | TCP | HTTP 访问 |
| 443 | TCP | HTTPS 访问 |

---

## 4. 数据安全

### 敏感配置管理

- 所有敏感配置（数据库密码、JWT 密钥、SMTP 授权码）通过环境变量或配置文件注入
- 配置文件已加入 `.gitignore`，不会提交到版本控制
- 示例配置文件仅包含占位符，不含真实凭据

### 数据库安全

- MySQL 仅允许本地连接（`127.0.0.1`）
- 数据库密码通过环境变量 `MYSQL_ROOT_PASSWORD` 设置
- 数据库表自动迁移，无需手动执行 SQL

---

## 5. HTTPS 加密

项目使用 Let's Encrypt 免费证书实现 HTTPS：

- Certbot 自动申请和续期证书
- HTTP 请求自动重定向到 HTTPS
- 证书续期由 systemd timer 自动管理

---

## 6. 部署安全

### .gitignore 配置

以下敏感文件已配置为不提交到版本控制：

```
config/.env
config/backend/.env
config/backend/config.dev.yaml
config/backend/config.prod.yaml
.env
```

### 容器镜像安全

- 基础镜像使用官方版本（`mysql:8.0`、`redis:7-alpine`）
- 定期更新基础镜像以获取安全补丁
