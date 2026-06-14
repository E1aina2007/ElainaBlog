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

#### 配置 Docker 镜像加速（国内服务器必做）

MySQL、Redis 等官方镜像在 `docker compose up` 时自动拉取，国内服务器需配置镜像加速，否则可能超时失败：

```bash
sudo tee /etc/docker/daemon.json > /dev/null << 'EOF'
{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://docker.xuanyuan.me"
  ]
}
EOF

sudo systemctl daemon-reload
sudo systemctl restart docker
```

#### 配置 Swap 空间（≤2GB 内存服务器必做）

Docker 构建镜像时内存占用较高，Swap 可作为安全垫防止 OOM 导致卡死：

```bash
# 创建 2GB Swap 文件
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# 持久化（重启后自动挂载）
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

# 验证
free -h
```

> Swap 仅作为构建时的安全垫，运行时不会被使用（Docker 容器默认不使用 Swap）。

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
  config: charset=utf8mb4&parseTime=True  # 不要添加 loc=Local，时区已由代码自动设置为 Asia/Shanghai

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

### 3. 初始化目录

```bash
# 创建绑定挂载的宿主机目录
mkdir -p config/backend
mkdir -p backend/uploads
mkdir -p frontend/public/author
```

#### 挂载目录一览

| 宿主机路径 | 容器内路径 | 服务 | 用途 |
|-----------|-----------|------|------|
| `config/backend` | `/app/config` | backend | 后端配置文件（`.env` 和 `yaml`） |
| `backend/uploads` | `/app/uploads` | backend | 用户上传文件（头像、文章封面等） |
| `frontend/public/author` | `/usr/share/nginx/html/author` | frontend | 作者头像与背景图 |
| `logs` (Volume) | `/app/log` | backend | 后端运行日志 |
| `mysql_data` (Volume) | `/var/lib/mysql` | mysql | MySQL 数据库数据 |
| `redis_data` (Volume) | `/data` | redis | Redis 缓存数据 |

> Docker 命名卷（`logs`、`mysql_data`、`redis_data`）会在首次 `docker compose up` 时自动创建。

#### 作者资源

作者头像和背景图存放在 `frontend/public/author/` 目录，挂载到前端容器内供 Nginx 直接提供静态访问。修改后**无需重启容器**。

```
frontend/public/author/
├── avatar.jpg      # 作者头像
└── background.jpg  # 作者页顶部背景图
```

### 4. 构建镜像

项目通过不同的 Compose 覆盖文件支持三种构建方式，按需选择：

| 方式 | 适用场景 | 启动命令 |
|------|---------|---------|
| 远程镜像 | 服务器无需构建，直接拉取 | `docker compose up -d` |
| 本地构建 | 服务器内存 ≥2GB | `docker compose -f docker-compose.yml -f docker-compose.local.yml up -d --build` |
| 本地编译上传 | 服务器内存 <2GB | `docker compose -f docker-compose.yml -f docker-compose.cross.yml up -d` |

> MySQL 和 Redis 使用官方镜像（`mysql:8.0`、`redis:7-alpine`），无需构建，会在 `docker compose up` 时自动拉取。如未提前配置镜像加速，请先完成「环境准备」中的镜像加速配置。

#### 方式一：远程镜像（推荐，适合 CI/CD 流程）

GitHub 推送到 main 分支时自动构建镜像并推送到远程仓库，服务器直接拉取。

**前置条件：** 已配置 GitHub Actions + 远程容器仓库（详见 [GITHUB_ACTIONS_DOCKER.md](docs/GITHUB_ACTIONS_DOCKER.md)）。

在 `.env` 中配置镜像地址：

```env
DOCKER_REGISTRY=registry.<地域>.aliyuncs.com
DOCKER_NAMESPACE=<命名空间>
```

```bash
# 拉取远程镜像并启动
docker compose pull
docker compose up -d
```

#### 方式二：服务器本地构建（适合内存 ≥2GB 的服务器）

```bash
# 顺序构建前后端镜像并启动（--parallel=false 防止同时构建导致内存不足）
docker compose -f docker-compose.yml -f docker-compose.local.yml up -d --build
```

> ⚠️ 不要省略 `--parallel=false`。前后端同时构建可能因内存不足（OOM）导致构建卡死或失败，尤其在 ≤2GB 内存的服务器上。

```bash
# 确认镜像构建成功
docker images | grep elainablog
```

预期输出：

```
elainablog-frontend   latest   ...   ...
elainablog-backend    latest   ...   ...
```

#### 方式三：本地编译后上传（适合内存 <2GB 的服务器）

服务器内存不足时，可在本地编译前后端产物，上传到服务器后构建镜像，跳过服务器上的编译过程。

**Step 1：本地编译**

后端交叉编译（在本地开发机器上执行）：

```bash
cd backend

# Windows (PowerShell)
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o elainablog ./cmd

# macOS / Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o elainablog ./cmd
```

前端构建：

```bash
cd frontend
npm run build
```

**Step 2：上传文件到服务器**

```bash
# 上传后端二进制和迁移脚本
scp backend/elainablog user@your-server-ip:~/ElainaBlog/backend/
scp -r backend/config/db/SQLscript user@your-server-ip:~/ElainaBlog/backend/config/db/

# 上传前端构建产物和 nginx 配置
scp -r frontend/dist user@your-server-ip:~/ElainaBlog/frontend/
scp frontend/nginx.conf user@your-server-ip:~/ElainaBlog/frontend/
```

**Step 3：构建并启动**

```bash
cd ~/ElainaBlog
docker compose -f docker-compose.yml -f docker-compose.cross.yml up -d
```

> 💡 更新时只需重复 Step 1 → Step 2 → Step 3，无需重新克隆项目。

### 5. 启动数据库

```bash
# 先单独启动 MySQL，等待就绪
docker compose up -d mysql
```

```bash
# 查看 MySQL 日志，出现 "ready for connections" 表示就绪
docker compose logs -f mysql
# 按 Ctrl+C 退出日志查看
```

```bash
# 确认 MySQL 容器状态为 running (healthy)
docker compose ps mysql
```

预期输出：

```
NAME                STATUS
elainablog-mysql    running (healthy)
```

### 6. 启动 Redis

```bash
# 启动 Redis
docker compose up -d redis
```

```bash
# 确认 Redis 容器状态为 running (healthy)
docker compose ps redis
```

预期输出：

```
NAME                STATUS
elainablog-redis    running (healthy)
```

### 7. 启动后端

```bash
# 启动后端容器（会自动执行数据库迁移）
docker compose up -d backend
```

```bash
# 查看后端启动日志，确认无报错
docker compose logs backend
```

```bash
# 确认后端容器状态为 running
docker compose ps backend
```

```bash
# 测试后端健康检查
curl http://127.0.0.1:9178/health
```

预期输出：

```
NAME                STATUS
elainablog-backend  running
```

> 数据库表会在后端首次启动时自动迁移创建，无需手动执行 SQL 脚本。

### 8. 启动前端

```bash
# 启动前端容器
docker compose up -d frontend
```

```bash
# 确认前端容器状态为 running
docker compose ps frontend
```

```bash
# 测试前端是否可访问
curl -I http://127.0.0.1:3000
```

预期输出：

```
NAME                 STATUS
elainablog-frontend  running
```

### 9. 确认全部服务

```bash
# 确认所有 4 个容器均为 running
docker compose ps
```

预期输出：

```
NAME                 STATUS
elainablog-frontend  running
elainablog-backend   running
elainablog-mysql     running (healthy)
elainablog-redis     running (healthy)
```

### 10. 创建管理员账号

```bash
# 方式一：命令行传入密码（推荐）
docker exec elainablog-backend ./elainablog initSystem <管理员密码>
```

```bash
# 方式二：密码已在 config.prod.yaml 的 admin.password 中配置
docker exec elainablog-backend ./elainablog initSystem
```

执行成功后会输出管理员邮箱和初始化结果。

### 11. 迁移头像文件命名（可选）

如果从旧版本升级，头像文件可能使用邮箱作为文件名。执行以下命令将其迁移为哈希命名：

```bash
docker exec elainablog-backend ./elainablog migrateAvatars
```

执行成功后会输出迁移结果。新上传的头像会自动使用哈希命名，无需迁移。

### 12. 配置宿主机 Nginx

```bash
# 创建 Nginx 配置文件
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
# 检查 Nginx 配置语法
sudo nginx -t
```

预期输出：

```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

```bash
# 重载 Nginx 配置
sudo nginx -s reload
```

> 如需 HTTPS，可使用以下命令自动申请 Let's Encrypt 证书并配置：
>
> ```bash
> sudo apt install -y certbot python3-certbot-nginx
> sudo certbot --nginx -d your-domain.com -d www.your-domain.com
> ```

### 13. 验证服务

```bash
# 测试后端健康检查
curl http://127.0.0.1:9178/health
```

```bash
# 测试前端
curl -I http://127.0.0.1:3000
```

```bash
# 测试宿主机 Nginx 反向代理
curl -I http://localhost
```

```bash
# 测试域名访问（DNS 解析生效后）
curl -I http://your-domain.com
curl http://your-domain.com/api/ui/health
```

---

## 三、更新部署

```bash
# 进入项目目录
cd ElainaBlog
```

```bash
# 拉取最新代码
git pull
```

#### 远程镜像模式

```bash
docker compose pull
docker compose up -d
```

#### 本地构建模式

```bash
docker compose -f docker-compose.yml -f docker-compose.local.yml up -d --build
```

#### 交叉编译模式

```bash
# 上传新编译的二进制后
scp backend/elainablog server:~/ElainaBlog/backend/
docker compose -f docker-compose.yml -f docker-compose.cross.yml up -d
```

```bash
# 确认所有容器均为 running
docker compose ps
```

```bash
# 查看后端日志，确认启动正常
docker compose logs backend
```

#### 仅更新前端（本地构建模式）

```bash
docker compose -f docker-compose.yml -f docker-compose.local.yml up -d --build frontend
```

#### 仅更新后端（本地构建模式）

```bash
docker compose -f docker-compose.yml -f docker-compose.local.yml up -d --build backend
docker compose build backend
docker compose up -d backend
docker compose logs backend
```

---

## 四、运维命令速查

### 容器管理

```bash
# 查看所有容器状态
docker compose ps
```

```bash
# 启动所有服务
docker compose up -d
```

```bash
# 停止所有服务（保留数据卷）
docker compose down
```

```bash
# 停止并删除数据卷（⚠️ 会删除数据库、日志、Redis 缓存）
docker compose down -v
```

```bash
# 重启所有服务
docker compose restart
```

```bash
# 重启单个服务
docker compose restart backend
docker compose restart frontend
docker compose restart mysql
docker compose restart redis
```

### 日志查看

```bash
# 查看后端日志（实时，Ctrl+C 退出）
docker compose logs -f backend
```

```bash
# 查看前端日志（实时）
docker compose logs -f frontend
```

```bash
# 查看 MySQL 日志（实时）
docker compose logs -f mysql
```

```bash
# 查看 Redis 日志（实时）
docker compose logs -f redis
```

```bash
# 查看所有服务最近 100 行日志
docker compose logs --tail 100
```

### 容器操作

```bash
# 进入后端容器
docker exec -it elainablog-backend sh
```

```bash
# 进入 MySQL 容器
docker exec -it elainablog-mysql mysql -u root -p
```

```bash
# 进入 Redis 容器
docker exec -it elainablog-redis redis-cli
```

```bash
# 查看后端容器内文件
docker exec elainablog-backend ls /app/uploads
docker exec elainablog-backend ls /app/log
docker exec elainablog-backend ls /app/migrations
```

### 数据备份与恢复

```bash
# 备份 MySQL 数据库
docker exec elainablog-mysql mysqldump -u root -p ElainaBlog > backup_$(date +%Y%m%d).sql
```

```bash
# 恢复 MySQL 数据库
docker exec -i elainablog-mysql mysql -u root -p ElainaBlog < backup.sql
```

```bash
# 备份 Redis
docker exec elainablog-redis redis-cli BGSAVE
docker cp elainablog-redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d).rdb
```

```bash
# 备份上传文件
cp -r backend/uploads backend/uploads_backup_$(date +%Y%m%d)
```

### 资源监控

```bash
# 查看容器资源占用（CPU、内存、网络）
docker stats --no-stream
```

```bash
# 查看磁盘使用
docker system df
```

```bash
# 清理无用镜像和缓存
docker system prune -f
```

---

## 五、常见问题

### 后端启动失败，提示数据库连接超时

```bash
# 确认 MySQL 容器已就绪
docker compose logs mysql
```

等待出现 "ready for connections" 后再启动后端。

### 前端页面空白或 404

```bash
# 1. 确认前端容器运行中
docker compose ps

# 2. 确认前端可访问
curl http://127.0.0.1:3000

# 3. 检查宿主机 Nginx 配置是否正确代理到 127.0.0.1:3000
cat /etc/nginx/conf.d/elainablog.conf
```

### API 请求 404 或 502

```bash
# 1. 确认后端容器运行中
docker compose ps

# 2. 确认端口绑定
curl http://127.0.0.1:9178/health

# 3. 查看后端日志排查错误
docker compose logs backend
```

### 文件上传失败

```bash
# 确认 uploads 目录已正确挂载
docker exec elainablog-backend ls /app/uploads

# 确认容器内目录权限
docker exec elainablog-backend ls -la /app/uploads
```

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
