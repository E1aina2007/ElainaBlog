# ElainaBlog 部署指南

---

## 目录

- [环境配置](#环境配置)
- [一、本地开发](#一本地开发)
- [二、生产环境部署](#二生产环境部署)
- [三、更新部署](#三更新部署)
- [四、运维命令速查](#四运维命令速查)
- [五、常见问题](#五常见问题)
- [六、上线检查清单](#六上线检查清单)

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

## 二、生产环境部署

> 架构：宿主机 Nginx 作为入口反向代理 → Docker Compose 运行前端（Nginx 托管静态文件）+ 后端（Go + MySQL + Redis），前后端均容器化。

### 环境准备

部署前确保服务器满足以下条件：

- 操作系统：Ubuntu 22.04 LTS / Debian 12
- 已安装 Docker 和 Docker Compose
- 已安装 Nginx
- 防火墙开放 22/80/443 端口
- 如使用国内服务器，已完成 ICP 备案
- 域名 DNS 已解析到服务器 IP

### 1. 克隆项目

```bash
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
mkdir -p config/backend backend/uploads frontend/public/author
```

> Docker 命名卷（`logs`、`mysql_data`、`redis_data`）会在首次 `docker compose up` 时自动创建。

#### 作者资源

作者头像和背景图存放在 `frontend/public/author/` 目录，挂载到前端容器内供 Nginx 直接提供静态访问。修改后**无需重启容器**。

```
frontend/public/author/
├── avatar.jpg      # 作者头像
└── background.jpg  # 作者页顶部背景图
```

#### 数据备份

```bash
# 备份 MySQL
docker exec elainablog-mysql mysqldump -u root -p ElainaBlog > backup_$(date +%Y%m%d).sql

# 恢复 MySQL
docker exec -i elainablog-mysql mysql -u root -p ElainaBlog < backup.sql

# 备份 Redis
docker exec elainablog-redis redis-cli BGSAVE
docker cp elainablog-redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d).rdb
```

> ⚠️ `docker compose down -v` 会删除所有数据卷（数据库、日志、Redis 缓存将全部丢失）。仅需重启服务请使用 `docker compose down && docker compose up -d`。

### 4. 启动服务

```bash
# 构建镜像并启动全部容器（frontend、backend、mysql、redis）
docker compose up -d --build

# 确认所有容器均为 running 状态
docker compose ps
```

预期输出应包含 4 个容器，状态均为 `running`：

```
NAME                STATUS
elainablog-frontend running
elainablog-backend  running
elainablog-mysql    running (healthy)
elainablog-redis    running (healthy)
```

```bash
# 查看后端启动日志，确认无报错
docker compose logs backend

# 测试后端健康检查
curl http://127.0.0.1:9178/health
```

> 数据库表会在后端首次启动时自动迁移创建，无需手动执行 SQL 脚本。

### 5. 创建管理员账号

```bash
# 方式一：命令行传入密码（推荐）
docker exec elainablog-backend ./elainablog initSystem <管理员密码>

# 方式二：密码已在 config.prod.yaml 的 admin.password 中配置
docker exec elainablog-backend ./elainablog initSystem
```

### 6. 配置宿主机 Nginx

```bash
sudo tee /etc/nginx/conf.d/elainablog.conf > /dev/null << 'EOF'
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
EOF
```

```bash
# 检查语法并重载
sudo nginx -t && sudo nginx -s reload
```

> 如需 HTTPS，可使用 `sudo certbot --nginx -d your-domain.com` 自动申请 Let's Encrypt 证书并配置。

### 7. 验证服务

```bash
curl http://127.0.0.1:9178/health    # 后端
curl -I http://127.0.0.1:3000        # 前端
curl -I http://your-domain.com       # 域名访问
```

---

## 三、更新部署

```bash
cd ElainaBlog
git pull

# 重新构建并重启所有服务
docker compose build
docker compose up -d

# 如只需更新前端或后端
docker compose build frontend && docker compose up -d frontend
docker compose build backend && docker compose up -d backend
```

---

## 四、运维命令速查

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

# 查看系统资源
docker stats --no-stream
```

---

## 五、常见问题

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
docker exec elainablog-backend ./elainablog initSystem <新密码>
```

---

## 六、上线检查清单

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
- [ ] HTTPS 证书已配置并可访问（如需要）

### 安全

- [ ] 防火墙已启用（仅开放 22/80/443）
- [ ] 数据库端口未暴露到公网
