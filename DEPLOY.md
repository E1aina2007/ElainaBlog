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
  - [4.1 Docker 挂载目录说明](#41-docker-挂载目录说明)
- [五、HTTPS 配置](#五https-配置)
- [六、Cloudflare 配置](#六cloudflare-配置)
- [七、更新部署](#七更新部署)
- [八、运维命令速查](#八运维命令速查)
- [九、常见问题](#九常见问题)
- [十、上线检查清单](#十上线检查清单)

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

在域名注册商的 DNS 管理面板中添加 A 记录，将域名指向服务器公网 IP：

| 记录类型 | 主机记录 | 记录值 | TTL |
|---------|---------|-------|-----|
| A | `@` | 服务器公网 IP | 600 |
| A | `www` | 服务器公网 IP | 600 |

#### 各平台操作步骤

**阿里云：**

1. 登录 [阿里云控制台](https://dns.console.aliyun.com)
2. 进入「域名解析 DNS」→「解析设置」
3. 点击「添加记录」，填写主机记录和记录值

**腾讯云：**

1. 登录 [DNSPod 控制台](https://console.dnspod.cn)
2. 选择域名 →「添加记录」
3. 填写主机记录、记录类型和记录值

**Cloudflare（推荐海外服务器）：**

1. 登录 [Cloudflare](https://dash.cloudflare.com)
2. 添加站点并按提示修改域名的 NS 服务器
3. 在「DNS」→「Records」中添加 A 记录
4. 代理状态选择「已代理」（橙色云朵）可启用 CDN 和防护

#### 验证 DNS 解析

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
# 选择一个目录存放项目，如 /root、/opt、/home/user 等
cd ~
git clone <repository-url> ElainaBlog
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

> 生成随机密钥（用于 `access_token_secret`、`refresh_token_secret`、`sessions_key`）：
>
> ```bash
> # Linux / macOS
> openssl rand -hex 32
>
> # PowerShell
> -join ((1..32) | ForEach-Object { '{0:x2}' -f (Get-Random -Max 256) })
> ```
>
> 三个密钥应分别生成，不要相同。

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

### 3. Docker 挂载目录说明

项目使用两种挂载方式：**绑定挂载**（Bind Mount）将宿主机目录直接映射到容器，**命名卷**（Named Volume）由 Docker 自动管理。

#### 挂载目录一览

| 宿主机路径 | 容器内路径 | 服务 | 用途 |
|-----------|-----------|------|------|
| `config/backend` | `/app/config` | backend | 后端配置文件（`.env` 和 `yaml`） |
| `backend/uploads` | `/app/uploads` | backend | 用户上传文件（头像、文章封面等） |
| `frontend/public/author` | `/usr/share/nginx/html/author` | frontend | 作者头像与背景图 |
| `logs` (Volume) | `/app/log` | backend | 后端运行日志 |
| `mysql_data` (Volume) | `/var/lib/mysql` | mysql | MySQL 数据库数据 |
| `redis_data` (Volume) | `/data` | redis | Redis 缓存数据 |

#### 目录初始化

首次部署时，需确保绑定挂载的宿主机目录存在：

```bash
# 创建配置目录（如尚未创建）
mkdir -p config/backend

# 创建上传文件目录
mkdir -p backend/uploads

# 创建前端作者资源目录
mkdir -p frontend/public/author
```

> Docker 命名卷（`logs`、`mysql_data`、`redis_data`）会在首次 `docker compose up` 时自动创建，无需手动处理。

#### 配置文件挂载

后端配置文件通过 `config/backend` 目录挂载到容器内 `/app/config`，修改宿主机上的配置文件后重启后端容器即可生效：

```bash
# 编辑配置
vim config/backend/config.prod.yaml

# 重启后端使配置生效
docker compose restart backend
```

#### 作者资源管理

作者头像和背景图存放在 `frontend/public/author/` 目录，挂载到前端容器内 `/usr/share/nginx/html/author`，通过 Nginx 直接提供静态访问。

**目录结构：**

```
frontend/public/author/
├── avatar.jpg      # 作者头像
└── background.jpg  # 作者页顶部背景图
```

**上传到服务器的方式：**

文件必须放到宿主机的挂载目录 `frontend/public/author/`（项目根目录下），Docker 会自动将其映射到容器内供 Nginx 访问。

方式一：SCP 从本地传输（推荐，适合首次部署或少量文件更新）

```bash
# 从本地上传整个 author 目录到服务器的挂载目录
# 将 <项目路径> 替换为服务器上的实际路径，如 /root/ElainaBlog、/opt/ElainaBlog
scp -r frontend/public/author/* root@your-server:<项目路径>/frontend/public/author/
```

方式二：rsync 增量同步（推荐，适合频繁更新或大量文件）

```bash
# 增量同步，仅传输有变化的文件到服务器的挂载目录
rsync -avz --progress frontend/public/author/ root@your-server:<项目路径>/frontend/public/author/
```

方式三：通过 Git 同步（适合已纳入版本管理的资源）

```bash
# 在服务器上拉取最新代码（author 目录已在 public 下）
cd ElainaBlog  # 进入项目根目录
git pull
# 无需重启容器，Nginx 会直接读取挂载目录中的新文件
```

> 修改 `frontend/public/author/` 下的文件后**无需重启容器**，因为该目录通过绑定挂载映射到容器内，Nginx 会直接读取最新文件。

#### 上传文件管理

用户上传的头像、文章封面等文件存储在 `backend/uploads` 目录，挂载到容器内 `/app/uploads`。备份时可直接复制该目录：

```bash
# 备份上传文件
cp -r backend/uploads backend/uploads_backup_$(date +%Y%m%d)
```

#### 日志查看

后端日志通过命名卷 `logs` 持久化，可使用以下方式查看：

```bash
# 查看实时日志
docker compose logs -f backend

# 进入容器查看日志文件
docker exec -it elainablog-backend ls /app/log
```

#### 数据备份与恢复

**备份 MySQL：**

```bash
docker exec elainablog-mysql mysqldump -u root -p ElainaBlog > backup_$(date +%Y%m%d).sql
```

**恢复 MySQL：**

```bash
docker exec -i elainablog-mysql mysql -u root -p ElainaBlog < backup.sql
```

**备份 Redis（如需要）：**

```bash
docker exec elainablog-redis redis-cli BGSAVE
docker cp elainablog-redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d).rdb
```

#### 危险操作提醒

```bash
# ⚠️ 以下命令会删除所有数据卷（数据库、日志、Redis 缓存将全部丢失）
docker compose down -v

# 如仅需重启服务而不丢失数据，使用：
docker compose down
docker compose up -d
```

### 4. 启动数据库

```bash
docker compose up -d mysql

# 观察日志，出现 "ready for connections" 后继续
docker compose logs -f mysql
```

> 数据库表会在后端首次启动时自动迁移创建，无需手动执行 SQL 脚本。

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

### 8. 验证服务

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

## 六、Cloudflare 配置

Cloudflare 提供免费 CDN、DDoS 防护和 WAF 防火墙，推荐海外服务器使用。

### 1. 添加站点

1. 注册并登录 [Cloudflare](https://dash.cloudflare.com)
2. 点击「Add a site」，输入你的域名
3. 选择免费套餐（Free plan）
4. Cloudflare 会扫描现有 DNS 记录，确认无误后继续

### 2. 修改域名 NS 服务器

Cloudflare 会分配两个 NS 服务器，例如：

```
ns1.cloudflare.com
ns2.cloudflare.com
```

前往你的域名注册商控制台，将域名的 NS 服务器修改为 Cloudflare 提供的地址。NS 生效通常需要几分钟到 24 小时。

### 3. 配置 DNS 记录

在 Cloudflare「DNS」→「Records」中确认以下记录：

| 类型 | 名称 | 内容 | 代理状态 |
|------|------|------|---------|
| A | `@` | 服务器公网 IP | 已代理（橙色云朵） |
| A | `www` | 服务器公网 IP | 已代理（橙色云朵） |

> 代理状态为「已代理」时，流量经过 Cloudflare 隐藏源站 IP；「仅 DNS」则直接暴露 IP。

### 4. 配置 SSL/TLS

进入「SSL/TLS」→「Overview」：

- 加密模式选择 **Full (Strict)**（需要源站已配置 HTTPS 证书）
- 开启「Always Use HTTPS」（强制 HTTPS 重定向）
- 开启「Automatic HTTPS Rewrites」

### 5. 配置安全防护

进入「Security」→「Settings」：

| 设置项 | 推荐值 | 说明 |
|--------|--------|------|
| Security Level | Medium | 对可疑 IP 发起质询 |
| Browser Integrity Check | 开启 | 检测恶意请求头 |
| Challenge Passage | 30 minutes | 质询通过后的有效期 |

### 6. 配置缓存

进入「Caching」→「Configuration」：

- Browser Cache TTL：选择「Respect Existing Headers」
- 开发调试时可使用「Purge Cache」清除缓存

### 7. 验证 Cloudflare 生效

```bash
# 检查响应头是否包含 Cloudflare 标识
curl -I https://your-domain.com

# 应看到类似以下 header
# server: cloudflare
# cf-ray: xxxxxxxxx
```

### 8. 获取真实客户端 IP

由于流量经过 Cloudflare 代理，后端获取的客户端 IP 为 Cloudflare IP。如需获取真实 IP，需配置 Nginx 使用 `CF-Connecting-IP` 头：

```nginx
# 在 http 块中添加 Cloudflare IP 段
set_real_ip_from 173.245.48.0/20;
set_real_ip_from 103.21.244.0/22;
set_real_ip_from 103.22.200.0/22;
set_real_ip_from 103.31.4.0/22;
set_real_ip_from 141.101.64.0/18;
set_real_ip_from 108.162.192.0/18;
set_real_ip_from 190.93.240.0/20;
set_real_ip_from 188.114.96.0/20;
set_real_ip_from 197.234.240.0/22;
set_real_ip_from 198.41.128.0/17;
set_real_ip_from 162.158.0.0/15;
set_real_ip_from 104.16.0.0/13;
set_real_ip_from 104.24.0.0/14;
set_real_ip_from 172.64.0.0/13;
set_real_ip_from 131.0.72.0/22;
real_ip_header CF-Connecting-IP;
```

> Cloudflare IP 段可能会更新，完整列表见 https://www.cloudflare.com/ips/

---

## 七、更新部署

```bash
cd ElainaBlog  # 进入项目根目录
git pull

# 重新构建并重启所有服务
docker compose build
docker compose up -d

# 如只需更新前端或后端
docker compose build frontend && docker compose up -d frontend
docker compose build backend && docker compose up -d backend
```

---

## 八、运维命令速查

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

## 九、常见问题

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

## 十、上线检查清单

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
- [ ] Cloudflare 已配置并生效（如使用）
