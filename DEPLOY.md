# ElainaBlog 部署指南

本指南覆盖从零开始到网站实际上线的完整流程，包括本地开发、服务器准备、生产部署、安全加固、自动化运维等各环节。

---

## 目录

- [环境配置](#环境配置)
- [一、本地开发](#一本地开发)
- [二、上线前准备](#二上线前准备)
  - [2.1 服务器选购](#21-服务器选购)
  - [2.2 域名注册与备案](#22-域名注册与备案)
  - [2.3 域名解析](#23-域名解析)
- [三、服务器初始化](#三服务器初始化)
  - [3.1 系统基础配置](#31-系统基础配置)
  - [3.2 安装 Docker 与 Nginx](#32-安装-docker-与-nginx)
  - [3.3 安全加固](#33-安全加固)
- [四、生产环境部署](#四生产环境部署)
- [五、HTTPS 与域名接入](#五https-与域名接入)
- [六、安全加固（进阶）](#六安全加固进阶)
- [七、自动化备份](#七自动化备份)
- [八、监控与告警](#八监控与告警)
- [九、CI/CD 自动部署](#九cicd-自动部署)
- [十、SEO 与性能优化](#十seo-与性能优化)
- [十一、更新部署](#十一更新部署)
- [十二、运维命令速查](#十二运维命令速查)
- [十三、常见问题](#十三常见问题)
- [十四、上线检查清单](#十四上线检查清单)

---

## 环境配置

项目通过 `config/backend/.env` 中的 `MODE` 变量区分环境，加载对应的配置文件：

| MODE 值 | 加载的配置文件 | 用途 |
|---------|--------------|------|
| `dev` | `config/backend/config.dev.yaml` | 本地开发 |
| `prod` | `config/backend/config.prod.yaml` | 生产部署 |

所有配置文件统一放在 `config/backend/` 目录下，已加入 `.gitignore`，需从示例文件复制后修改：

```bash
cp config/backend/.env.example config/backend/.env
cp config/backend/config.example.yaml config/backend/config.dev.yaml
```

### 必须修改的配置项

| 配置项 | 说明 |
|-------|------|
| `auth.access_token_secret` | JWT Access Token 签名密钥 |
| `auth.refresh_token_secret` | JWT Refresh Token 签名密钥 |
| `db.password` | MySQL 密码 |
| `smtp.from` / `smtp.verification` | 发件邮箱及 SMTP 授权码 |
| `server.sessions_key` | Session 密钥 |

### 开发与生产配置差异

| 配置项 | 开发环境 | 生产环境 |
|-------|---------|---------|
| `db.host` | `127.0.0.1` | `mysql`（Docker 服务名） |
| `redis.address` | `127.0.0.1:6379` | `redis:6379`（Docker 服务名） |
| `server.host` | `127.0.0.1` | `0.0.0.0` |
| `server.env` | `debug` | `production` |
| `zap.is_console_print` | `true` | `false` |
| `dev` | `true` | `false` |

> 请勿将真实密码写入 `*.example.yaml` 示例文件。

---

## 一、本地开发

### 环境要求

| 组件 | 版本要求 |
|------|---------|
| Go | 1.25+ |
| Node.js | 20+ |
| MySQL | 8.0+ |
| Redis | 6.0+ |

### 1. 克隆项目

```bash
git clone <repository-url>
cd ElainaBlog
```

### 2. 配置环境

```bash
cp config/backend/.env.example config/backend/.env
cp config/backend/config.example.yaml config/backend/config.dev.yaml
```

编辑 `config/backend/config.dev.yaml`，修改数据库、Redis 等连接信息。

确保 `config/backend/.env` 中设置了运行模式：

```
MODE=dev
```

### 3. 初始化管理员

数据库表会在后端首次启动时自动迁移创建，无需手动执行 SQL。

```bash
cd backend
go run ./cmd initSystem <password>
```

### 4. 启动后端

```bash
cd backend
go run ./cmd runServer
```

后端默认运行在 `http://localhost:8080`

### 5. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端默认运行在 `http://localhost:5173`

---

## 二、上线前准备

### 2.1 服务器选购

推荐云服务器配置（以日均 PV < 5000 的个人博客为例）：

| 配置项 | 最低要求 | 推荐配置 |
|-------|---------|---------|
| CPU | 1 核 | 2 核 |
| 内存 | 1 GB | 2 GB |
| 硬盘 | 20 GB SSD | 40 GB SSD |
| 带宽 | 1 Mbps | 3~5 Mbps |
| 操作系统 | Ubuntu 22.04 LTS / Debian 12 | Ubuntu 22.04 LTS |

**主流云服务商参考：**

| 服务商 | 优势 |
|-------|------|
| 阿里云 | 国内生态最完善，备案流程成熟 |
| 腾讯云 | 性价比高，学生优惠力度大 |
| 华为云 | 政企场景支持好 |
| 雨云 / 搬瓦工 | 海外免备案（适合不想备案的场景） |

> 如果选择海外服务器则无需 ICP 备案，但国内访问速度可能较慢，建议搭配 CDN 使用。

### 2.2 域名注册与备案

#### 域名注册

在以下平台注册域名：

- 阿里云万网：https://wanwang.aliyun.com
- 腾讯云 DNSPod：https://dnspod.cloud.tencent.com
- Cloudflare（海外）：https://www.cloudflare.com

建议选择 `.com`、`.cn`、`.me` 等常见后缀。

#### ICP 备案（国内服务器必需）

> 使用中国大陆服务器必须完成 ICP 备案，否则域名无法解析到国内服务器。

**备案流程概览：**

1. 在云服务商控制台提交备案申请（阿里云/腾讯云均有引导式流程）
2. 准备材料：身份证正反面照片、域名证书、网站负责人核验照片
3. 云服务商初审（1~3 个工作日）
4. 工信部短信核验
5. 管局终审（7~20 个工作日，各省不同）
6. 备案通过后获取备案号，在网站底部展示

**备案号展示要求：**

在网站 Footer 中添加备案号链接，格式参考：

```html
<a href="https://beian.miit.gov.cn/" target="_blank">粤ICP备XXXXXXXX号</a>
```

> 备案期间可以先用 IP 地址直接访问进行测试。

### 2.3 域名解析

在域名注册商的 DNS 管理面板中添加以下记录：

| 记录类型 | 主机记录 | 记录值 | TTL |
|---------|---------|-------|-----|
| A | `@` | 服务器公网 IP | 600 |
| A | `www` | 服务器公网 IP | 600 |

> TTL 设置为 600 秒（10 分钟），方便上线初期调试；稳定运行后可改为 3600 或更高。

**验证 DNS 解析：**

```bash
# 在本地终端执行
nslookup your-domain.com
ping your-domain.com
```

确认解析到的 IP 是你的服务器公网 IP。

---

## 三、服务器初始化

### 3.1 系统基础配置

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 设置时区（确保日志和定时任务时间正确）
sudo timedatectl set-timezone Asia/Shanghai

# 安装常用工具
sudo apt install -y curl wget git vim unzip

# 创建部署用户（可选，推荐不用 root 直接操作）
sudo useradd -m -s /bin/bash deploy
sudo usermod -aG docker deploy
sudo usermod -aG sudo deploy
```

### 3.2 安装 Docker 与 Nginx

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
# 重新登录终端使 docker 组生效

# 验证
docker --version
docker compose version
```

```bash
# 安装 Nginx
sudo apt install -y nginx
sudo systemctl enable nginx
sudo systemctl start nginx
```

### 3.3 安全加固

#### 配置防火墙

```bash
# 使用 UFW（Ubuntu 默认防火墙）
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh          # SSH 端口
sudo ufw allow 80/tcp       # HTTP
sudo ufw allow 443/tcp      # HTTPS
sudo ufw enable
sudo ufw status verbose
```

#### SSH 安全加固

编辑 `/etc/ssh/sshd_config`：

```bash
sudo vim /etc/ssh/sshd_config
```

修改以下配置：

```text
PermitRootLogin no              # 禁止 root 登录
PasswordAuthentication no       # 禁用密码登录（需先配置密钥）
PubkeyAuthentication yes        # 启用密钥登录
Port 2222                       # 修改默认端口（可选）
MaxAuthTries 3                  # 最大尝试次数
```

配置密钥登录：

```bash
# 在本地机器生成密钥（如果还没有）
ssh-keygen -t ed25519 -C "your_email@example.com"

# 将公钥复制到服务器
ssh-copy-id -p 22 deploy@your-server-ip
```

重启 SSH 服务：

```bash
sudo systemctl restart sshd
```

#### 安装 fail2ban（防暴力破解）

```bash
sudo apt install -y fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

---

## 四、生产环境部署

> 架构：宿主机 Nginx 作为入口反向代理 → Docker Compose 运行前端（Nginx 托管静态文件）+ 后端（Go + MySQL + Redis），前后端均容器化。

### 1. 克隆项目

```bash
cd /opt
sudo git clone <repository-url> ElainaBlog
sudo chown -R $USER:$USER ElainaBlog
cd ElainaBlog
```

### 2. 配置环境

```bash
# Docker Compose 环境变量（项目根目录，供 docker-compose.yml 插值使用）
cp config/.env.example .env

# 后端配置（挂载到容器内）
cp config/backend/.env.example config/backend/.env
cp config/backend/config.example.yaml config/backend/config.prod.yaml
```

编辑 `.env`，设置数据库密码：

```env
MYSQL_ROOT_PASSWORD=替换为数据库密码
MYSQL_DATABASE=ElainaBlog
```

编辑 `config/backend/.env`，设置生产模式：

```
MODE=prod
```

编辑 `config/backend/config.prod.yaml`，以下为 **必须修改** 的配置项：

```yaml
auth:
  access_token_secret: "替换为随机密钥"   # openssl rand -hex 32
  refresh_token_secret: "替换为随机密钥"

db:
  host: mysql                            # Docker 服务名
  password: "替换为数据库密码"              # 必须与 .env 中 MYSQL_ROOT_PASSWORD 一致

smtp:
  from: "your_email@qq.com"
  verification: "your_smtp_auth_code"    # SMTP 授权码（非邮箱密码）

redis:
  address: redis:6379
  # password: "替换为Redis密码"           # 如设置了 Redis 密码则取消注释

server:
  host: 0.0.0.0
  env: production
  sessions_key: "替换为随机密钥"           # openssl rand -hex 32

zap:
  is_console_print: false

dev: false
```

### 3. 启动数据库

```bash
docker compose up -d mysql

# 观察日志，出现 "ready for connections" 后继续
docker compose logs -f mysql
```

> 数据库表会在后端首次启动时自动迁移创建，无需手动执行 SQL 脚本。

### 4. 准备作者资源（可选）

将作者头像和背景图放入 `frontend/public/author/` 目录（`avatar.jpg` 和 `background.jpg`）。

### 5. 启动所有服务

```bash
# 构建并启动
docker compose up -d

# 确认 4 个容器均为 running
docker compose ps
```

### 6. 创建管理员账号

```bash
# 方式一：命令行传入密码
docker exec elainablog-backend ./elainablog initSystem <管理员密码>

# 方式二：密码已在 config.prod.yaml 的 admin.password 中配置
docker exec elainablog-backend ./elainablog initSystem
```

### 7. 配置宿主机 Nginx

创建 `/etc/nginx/conf.d/elainablog.conf`：

```nginx
server {
    listen 80;
    server_name your-domain.com;  # 替换为你的域名

    client_max_body_size 20m;

    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # 公共反向代理配置
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # API 反向代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
    }

    location /uploads/ {
        proxy_pass http://127.0.0.1:8080;
    }

    location /health {
        proxy_pass http://127.0.0.1:8080;
    }

    # 前端
    location / {
        proxy_pass http://127.0.0.1:3000;
    }
}
```

检查配置并重载：

```bash
sudo nginx -t
sudo nginx -s reload
```

### 8. 验证服务

```bash
# 测试后端健康检查
curl http://127.0.0.1:8080/health

# 测试前端
curl -I http://127.0.0.1:3000

# 测试域名访问（DNS 解析生效后）
curl -I http://your-domain.com
curl http://your-domain.com/api/ui/
```

---

## 五、HTTPS 与域名接入

### 申请 SSL 证书（Let's Encrypt）

```bash
# 安装 Certbot
sudo apt install -y certbot python3-certbot-nginx

# 申请证书（Certbot 会自动修改 Nginx 配置）
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# 验证自动续期
sudo certbot renew --dry-run

# 确保自动续期定时任务已启用
sudo systemctl enable certbot.timer
sudo systemctl status certbot.timer
```

Certbot 会自动在 Nginx 配置中添加 HTTPS server block 并配置 HTTP → HTTPS 重定向。

### 手动 Nginx HTTPS 配置（如需自定义）

```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # ... 其余配置同上 ...
}
```

### 强制 HTTPS 重定向

确保所有 HTTP 请求跳转到 HTTPS：

```bash
# 测试重定向
curl -I http://your-domain.com
# 应返回 301 跳转到 https://
```

---

## 六、安全加固（进阶）

### 6.1 Nginx 安全配置

在 `nginx.conf` 的 `http` 块中添加：

```nginx
# 隐藏 Nginx 版本号
server_tokens off;

# 限制请求方法
if ($request_method !~ ^(GET|HEAD|POST|PUT|DELETE|PATCH)$) {
    return 405;
}

# 防止点击劫持
add_header X-Frame-Options "SAMEORIGIN" always;

# 防止 MIME 类型嗅探
add_header X-Content-Type-Options "nosniff" always;

# XSS 防护
add_header X-XSS-Protection "1; mode=block" always;
```

### 6.2 Docker 安全

- MySQL 和 Redis 端口仅绑定 `127.0.0.1`（`docker-compose.yml` 中已配置）
- 不要将 MySQL 3306 或 Redis 6379 端口暴露到公网
- 定期更新基础镜像：

```bash
docker compose pull
docker compose up -d
```

### 6.3 配置 Nginx 速率限制

在 `/etc/nginx/nginx.conf` 的 `http` 块中添加：

```nginx
# API 速率限制：每秒 10 个请求
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

# 登录接口更严格：每秒 2 个请求
limit_req_zone $binary_remote_addr zone=login_limit:10m rate=2r/s;
```

在对应的 `location` 块中应用：

```nginx
location /api/ {
    limit_req zone=api_limit burst=20 nodelay;
    proxy_pass http://127.0.0.1:8080;
}

location /api/auth/login {
    limit_req zone=login_limit burst=5 nodelay;
    proxy_pass http://127.0.0.1:8080;
}
```

### 6.4 自动安全更新

```bash
# 安装 unattended-upgrades
sudo apt install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

---

## 七、自动化备份

### 7.1 数据库自动备份脚本

创建 `/opt/scripts/backup-db.sh`：

```bash
#!/bin/bash
set -e

BACKUP_DIR="/opt/backups/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
KEEP_DAYS=7

mkdir -p "$BACKUP_DIR"

# 备份数据库
docker exec elainablog-mysql mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" ElainaBlog \
  | gzip > "$BACKUP_DIR/elainablog_${DATE}.sql.gz"

# 删除超过 N 天的备份
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +$KEEP_DAYS -delete

echo "[$(date)] 数据库备份完成: elainablog_${DATE}.sql.gz"
```

> 将 `MYSQL_ROOT_PASSWORD` 替换为实际密码，或改为从 `.env` 文件读取。

### 7.2 上传文件备份

创建 `/opt/scripts/backup-uploads.sh`：

```bash
#!/bin/bash
set -e

BACKUP_DIR="/opt/backups/uploads"
DATE=$(date +%Y%m%d_%H%M%S)
KEEP_DAYS=30

mkdir -p "$BACKUP_DIR"

tar -czf "$BACKUP_DIR/uploads_${DATE}.tar.gz" -C /opt/ElainaBlog/backend uploads/

find "$BACKUP_DIR" -name "*.tar.gz" -mtime +$KEEP_DAYS -delete

echo "[$(date)] 上传文件备份完成: uploads_${DATE}.tar.gz"
```

### 7.3 配置定时任务

```bash
# 添加执行权限
chmod +x /opt/scripts/backup-db.sh
chmod +x /opt/scripts/backup-uploads.sh

# 编辑 crontab
crontab -e
```

添加以下定时任务：

```cron
# 每天凌晨 3 点备份数据库
0 3 * * * /opt/scripts/backup-db.sh >> /var/log/backup-db.log 2>&1

# 每周日凌晨 4 点备份上传文件
0 4 * * 0 /opt/scripts/backup-uploads.sh >> /var/log/backup-uploads.log 2>&1
```

### 7.4 备份验证

定期手动验证备份是否可用：

```bash
# 检查备份文件
ls -lh /opt/backups/mysql/

# 测试恢复（在开发环境或临时数据库中）
gunzip < /opt/backups/mysql/elainablog_XXXXXXXX_XXXXXX.sql.gz | \
  docker exec -i elainablog-mysql mysql -u root -p ElainaBlog
```

---

## 八、监控与告警

### 8.1 服务健康检查脚本

创建 `/opt/scripts/health-check.sh`：

```bash
#!/bin/bash

DOMAIN="your-domain.com"
WEBHOOK_URL=""  # 钉钉/企微/Telegram Webhook（留空则仅记录日志）

check_service() {
    local name=$1
    local url=$2
    local status

    status=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 5 --max-time 10 "$url")

    if [ "$status" != "200" ]; then
        local msg="[$(date)] ALERT: $name 异常 (HTTP $status)"
        echo "$msg"

        if [ -n "$WEBHOOK_URL" ]; then
            curl -s -X POST "$WEBHOOK_URL" \
              -H "Content-Type: application/json" \
              -d "{\"content\": \"$msg\"}"
        fi

        return 1
    fi

    echo "[$(date)] OK: $name (HTTP $status)"
    return 0
}

check_service "后端 API" "http://127.0.0.1:8080/health"
check_service "前端页面" "http://127.0.0.1:3000"
check_service "域名 HTTPS" "https://$DOMAIN"
```

配置定时检查：

```bash
chmod +x /opt/scripts/health-check.sh
crontab -e
# 每 5 分钟检查一次
*/5 * * * * /opt/scripts/health-check.sh >> /var/log/health-check.log 2>&1
```

### 8.2 日志管理

后端日志使用 Zap + lumberjack 实现了自动轮转，存放在 Docker 卷 `logs` 中。

```bash
# 查看后端日志
docker compose logs -f backend

# 查看 Nginx 访问日志
sudo tail -f /var/log/nginx/access.log

# 查看 Nginx 错误日志
sudo tail -f /var/log/nginx/error.log

# 查看 Docker 容器资源占用
docker stats --no-stream
```

### 8.3 磁盘空间监控

```bash
# 检查磁盘使用
df -h

# 检查 Docker 占用
docker system df

# 清理无用的 Docker 资源
docker system prune -f
```

---

## 九、CI/CD 自动部署

> 以 GitHub Actions 为例。项目根目录创建 `.github/workflows/deploy.yml`。

### 9.1 配置 GitHub Secrets

在 GitHub 仓库 Settings → Secrets and variables → Actions 中添加：

| Secret 名称 | 说明 |
|------------|------|
| `SERVER_HOST` | 服务器 IP 或域名 |
| `SERVER_USER` | SSH 用户名 |
| `SERVER_SSH_KEY` | SSH 私钥 |
| `SERVER_PORT` | SSH 端口（默认 22） |

### 9.2 GitHub Actions 配置

创建 `.github/workflows/deploy.yml`：

```yaml
name: Deploy to Production

on:
  push:
    branches: [main]
    paths:
      - 'backend/**'
      - 'frontend/**'
      - 'docker-compose.yml'
      - 'frontend/Dockerfile'
      - 'backend/Dockerfile'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to server
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          port: ${{ secrets.SERVER_PORT }}
          script: |
            cd /opt/ElainaBlog
            git pull origin main
            docker compose build
            docker compose up -d
            echo "部署完成: $(date)"
```

### 9.3 服务器端 SSH 密钥配置

```bash
# 在服务器上为部署用户添加 GitHub Actions 的公钥
mkdir -p ~/.ssh
echo "your-github-actions-public-key" >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

### 9.4 多环境部署（可选）

如有 staging 环境，可扩展为：

```
push dev branch → 部署到 staging
push main branch → 部署到 production
```

---

## 十、SEO 与性能优化

### 10.1 SEO 基础配置

#### robots.txt

在 `frontend/public/robots.txt` 中添加：

```text
User-agent: *
Allow: /
Disallow: /admin/
Disallow: /api/

Sitemap: https://your-domain.com/sitemap.xml
```

#### 网站 `<head>` Meta 标签

确保前端页面包含以下 meta 标签（在 `index.html` 或 Vue 组件中）：

```html
<meta name="description" content="你的博客描述">
<meta name="keywords" content="博客,技术,个人网站">
<meta property="og:title" content="ElainaBlog">
<meta property="og:description" content="你的博客描述">
<meta property="og:type" content="website">
<meta property="og:url" content="https://your-domain.com">
```

#### 备案号展示

在 Footer 组件中添加 ICP 备案号（国内服务器必需）。

### 10.2 性能优化

#### Gzip 压缩

前端容器 Nginx 和宿主机 Nginx 均已配置 Gzip（见 `nginx.conf`）。确认生效：

```bash
curl -H "Accept-Encoding: gzip" -I https://your-domain.com
# 响应头应包含 Content-Encoding: gzip
```

#### 静态资源缓存

前端构建产物（Vite 带 hash 文件名）已配置 30 天缓存。宿主机 Nginx 配置：

```nginx
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?|ttf|eot)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
    access_log off;
}
```

#### CDN 加速（可选）

如果服务器在海外或需要加速国内访问，可使用 CDN：

- 阿里云 CDN
- 腾讯云 CDN
- Cloudflare（免费套餐）

配置方式：将域名 DNS 指向 CDN 提供商的 CNAME，CDN 回源到你的服务器 IP。

---

## 十一、更新部署

```bash
cd /opt/ElainaBlog
git pull

# 重新构建并重启所有服务
docker compose build
docker compose up -d

# 如只需更新前端或后端
docker compose build frontend && docker compose up -d frontend
docker compose build backend && docker compose up -d backend
```

> 更新前建议先备份数据库：参见 [七、自动化备份](#七自动化备份)。

---

## 十二、运维命令速查

```bash
# 查看所有容器状态
docker compose ps

# 查看后端日志（实时）
docker compose logs -f backend

# 查看 MySQL 日志
docker compose logs -f mysql

# 重启后端
docker compose restart backend

# 重启所有服务
docker compose restart

# 停止所有服务
docker compose down

# 停止并删除数据卷（⚠️ 会删除数据库和上传文件）
docker compose down -v

# 进入后端容器
docker exec -it elainablog-backend sh

# 备份数据库
docker exec elainablog-mysql mysqldump -u root -p ElainaBlog > backup_$(date +%Y%m%d).sql

# 查看 Nginx 错误日志
sudo tail -f /var/log/nginx/error.log

# 查看系统资源
docker stats --no-stream

# 清理 Docker 无用资源
docker system prune -f
```

---

## 十三、常见问题

### 后端启动失败，提示数据库连接超时

确认 MySQL 容器已就绪：`docker compose logs mysql`，等待出现 "ready for connections" 后再启动后端。

### 前端页面空白或 404

1. 确认前端容器运行中：`docker compose ps`
2. 确认前端可访问：`curl http://127.0.0.1:3000`
3. 检查宿主机 Nginx 配置是否正确代理到 `127.0.0.1:3000`

### API 请求 404 或 502

1. 确认后端容器运行中：`docker compose ps`
2. 确认端口绑定：`curl http://127.0.0.1:8080/health`
3. 宿主机 Nginx 的 `/api/` 应代理到 `127.0.0.1:8080`

### 文件上传失败

1. 确认宿主机 Nginx 和前端容器 Nginx 中 `client_max_body_size` 均为 20m
2. 确认 `uploads` 数据卷已正确挂载：`docker compose exec backend ls /app/uploads`

### 作者头像或背景图不显示

1. 确认图片文件存在于 `frontend/public/author/` 目录
2. 确认 Docker 挂载正确：`docker compose exec frontend ls /usr/share/nginx/html/author/`

### 忘记管理员密码

```bash
# 命令行传入新密码
docker exec elainablog-backend ./elainablog initSystem <新密码>

# 或修改 config/backend/config.prod.yaml 中 admin.password 后重启容器
docker compose restart backend
docker exec elainablog-backend ./elainablog initSystem
```

如需重置已有账号的密码，需直接操作数据库。

### HTTPS 证书过期

```bash
# 手动续期
sudo certbot renew

# 检查自动续期状态
sudo systemctl status certbot.timer
```

### 磁盘空间不足

```bash
# 查看磁盘使用
df -h

# 清理 Docker 资源
docker system prune -f

# 清理旧日志
sudo journalctl --vacuum-time=7d

# 清理 APT 缓存
sudo apt clean
```

---

## 十四、上线检查清单

### 基础部署

- [ ] 服务器已购买并完成系统初始化
- [ ] 域名已注册
- [ ] ICP 备案已完成（国内服务器）
- [ ] DNS 解析已配置并生效
- [ ] Docker 和 Docker Compose 已安装
- [ ] Nginx 已安装并运行

### 配置正确性

- [ ] `config/backend/config.prod.yaml` 中所有密钥已替换为随机值
- [ ] `config/backend/config.prod.yaml` 中 `db.password` 与 `.env` 中 `MYSQL_ROOT_PASSWORD` 一致
- [ ] `config/backend/config.prod.yaml` 中 `dev: false`、`server.env: production`
- [ ] `config/backend/config.prod.yaml` 中 `admin` 段已配置管理员用户名和邮箱
- [ ] `config/backend/.env` 中 `MODE=prod`
- [ ] `.env` 中数据库密码已设置

### 服务运行

- [ ] `docker compose ps` 显示 4 个容器均为 running
- [ ] 管理员账号已创建
- [ ] 作者头像与背景图已放入 `frontend/public/author/`（可选）

### 域名与网络

- [ ] 宿主机 Nginx 配置中 `server_name` 已修改为实际域名
- [ ] Nginx 配置语法检查通过（`nginx -t`）
- [ ] `curl http://your-domain.com` 返回页面
- [ ] `curl http://your-domain.com/api/ui/` 返回 API 响应
- [ ] HTTPS 证书已配置并可访问
- [ ] HTTP → HTTPS 重定向已生效

### 安全

- [ ] 防火墙已配置（仅开放 22/80/443）
- [ ] SSH 已加固（禁用密码登录 / 修改默认端口）
- [ ] fail2ban 已安装并运行
- [ ] 数据库端口未暴露到公网
- [ ] Nginx 安全头已配置

### 运维

- [ ] 数据库自动备份已配置
- [ ] 上传文件备份已配置
- [ ] 健康检查脚本已配置
- [ ] 日志轮转正常工作

### SEO（可选）

- [ ] `robots.txt` 已配置
- [ ] Meta 标签已完善
- [ ] 备案号已在 Footer 展示（国内服务器）
- [ ] Sitemap 已生成（如需要）
