# ElainaBlog 部署指南

本指南覆盖从零开始到网站实际上线的完整流程，包括本地开发和生产部署的核心步骤。

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
  - [3.3 防火墙配置](#33-防火墙配置)
- [四、生产环境部署](#四生产环境部署)
- [五、HTTPS 配置](#五https-配置)
- [六、更新部署](#六更新部署)
- [七、运维命令速查](#七运维命令速查)
- [八、常见问题](#八常见问题)
- [九、上线检查清单](#九上线检查清单)

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

后端默认运行在 `http://localhost:9178`

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

> 海外服务器无需 ICP 备案，但国内访问速度可能较慢。

### 2.2 域名注册与备案

#### 域名注册

在以下平台注册域名：

- 阿里云万网：https://wanwang.aliyun.com
- 腾讯云 DNSPod：https://dnspod.cloud.tencent.com
- Cloudflare（海外）：https://www.cloudflare.com

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

在网站 Footer 中添加备案号链接：

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

**验证 DNS 解析：**

```bash
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

### 3.3 防火墙配置

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

### 4. 启动所有服务

```bash
# 构建并启动
docker compose up -d

# 确认 4 个容器均为 running
docker compose ps
```

### 5. 创建管理员账号

```bash
# 方式一：命令行传入密码
docker exec elainablog-backend ./elainablog initSystem <管理员密码>

# 方式二：密码已在 config.prod.yaml 的 admin.password 中配置
docker exec elainablog-backend ./elainablog initSystem
```

### 6. 配置宿主机 Nginx

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
        proxy_pass http://127.0.0.1:9178;
    }

    location /uploads/ {
        proxy_pass http://127.0.0.1:9178;
    }

    location /health {
        proxy_pass http://127.0.0.1:9178;
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

### 7. 验证服务

```bash
# 测试后端健康检查
curl http://127.0.0.1:9178/health

# 测试前端
curl -I http://127.0.0.1:3000

# 测试域名访问（DNS 解析生效后）
curl -I http://your-domain.com
curl http://your-domain.com/api/ui/
```

---

## 五、HTTPS 配置

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

### 验证 HTTPS

```bash
# 测试重定向
curl -I http://your-domain.com
# 应返回 301 跳转到 https://

# 测试 HTTPS
curl -I https://your-domain.com
```

---

## 六、更新部署

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

---

## 七、运维命令速查

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
```

---

## 八、常见问题

### 后端启动失败，提示数据库连接超时

确认 MySQL 容器已就绪：`docker compose logs mysql`，等待出现 "ready for connections" 后再启动后端。

### 前端页面空白或 404

1. 确认前端容器运行中：`docker compose ps`
2. 确认前端可访问：`curl http://127.0.0.1:3000`
3. 检查宿主机 Nginx 配置是否正确代理到 `127.0.0.1:3000`

### API 请求 404 或 502

1. 确认后端容器运行中：`docker compose ps`
2. 确认端口绑定：`curl http://127.0.0.1:9178/health`
3. 宿主机 Nginx 的 `/api/` 应代理到 `127.0.0.1:9178`

### 文件上传失败

1. 确认宿主机 Nginx 和前端容器 Nginx 中 `client_max_body_size` 均为 20m
2. 确认 `uploads` 数据卷已正确挂载：`docker compose exec backend ls /app/uploads`

### 忘记管理员密码

```bash
# 命令行传入新密码
docker exec elainablog-backend ./elainablog initSystem <新密码>
```

### HTTPS 证书过期

```bash
# 手动续期
sudo certbot renew

# 检查自动续期状态
sudo systemctl status certbot.timer
```

---

## 九、上线检查清单

### 配置

- [ ] `config/backend/config.prod.yaml` 中所有密钥已替换为随机值
- [ ] `config/backend/config.prod.yaml` 中 `db.password` 与 `.env` 中 `MYSQL_ROOT_PASSWORD` 一致
- [ ] `config/backend/config.prod.yaml` 中 `dev: false`、`server.env: production`
- [ ] `config/backend/.env` 中 `MODE=prod`
- [ ] `.env` 中数据库密码已设置

### 服务

- [ ] `docker compose ps` 显示 4 个容器均为 running
- [ ] `curl http://127.0.0.1:9178/health` 返回正常
- [ ] 管理员账号已创建

### 网络

- [ ] 域名 DNS 解析已指向服务器 IP
- [ ] 宿主机 Nginx 配置中 `server_name` 已修改为实际域名
- [ ] `nginx -t` 语法检查通过
- [ ] `curl http://your-domain.com` 返回页面
- [ ] HTTPS 证书已配置并可访问
- [ ] HTTP → HTTPS 重定向已生效

### 安全

- [ ] 防火墙已启用（仅开放 22/80/443）
- [ ] 数据库端口未暴露到公网
