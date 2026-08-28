# ElainaBlog

> 基于 Gin + Vue 3 二次开发的个人博客网站框架
>
> *Elaina 这个名字来源于动漫《魔女之旅》中的主角灰之魔女伊蕾娜。她很可爱。*

## 功能

### 文章系统

- Markdown 编辑与渲染，支持代码高亮（highlight.js）、表格、任务列表等扩展语法
- 文章分类管理，首页按分类筛选，文章置顶
- 草稿与发布状态切换

### 用户系统

- 邮箱注册 + 邮箱验证码验证
- JWT 双 Token 认证（Access Token + Refresh Token），无感刷新

### 评论系统

- 文章评论与回复
- 登录用户可发表、删除自己的评论
- 管理员可管理全部评论

### 作者主页

- 个人介绍、头像与背景图展示
- 博客统计数据（文章数、总字数、总阅读量、运行天数等）

### 管理后台（管理员）

- 仪表盘：博客核心数据概览
- 模块化管理：支持对用户，文章等模块内容的管理
- 系统状态：服务器运行状态监控
- 安全面板：IP 封禁管理
- 数据备份：数据库导出
- 缓存管理：一键清理缓存

### 在线工具

- 随机数生成器（支持范围、数量、去重）
- 时间戳转换（Unix 时间戳与日期互转）
- 颜色值转换（HEX / RGB / HSL，带预览色块）
- Markdown 在线预览（左右分栏，支持代码高亮）
- 正则表达式测试（实时匹配、高亮预览、语法速查表）

### 安全设计

**网络层**：Cloudflare CDN 代理隐藏源站 IP，WAF 规则拦截 SQL 注入 / XSS / 敏感文件访问，UFW 防火墙仅开放 22/80/443 端口。

**应用层**：
- GORM 参数化查询防 SQL 注入
- 密码 bcrypt 加盐哈希，JWT 双 Token 认证 + Redis 黑名单吊销
- Redis 滑动窗口限流（登录 / 注册 / 刷新 / 发码等敏感接口），IP 封禁机制
- 上传白名单 + MIME 魔数检测 + 随机文件名 + 大小限制
- 角色权限控制（普通用户 / 管理员），细分错误码体系（20+）

**服务端**：Nginx 安全头（`X-Frame-Options` / `X-Content-Type-Options` / `X-XSS-Protection` / `Referrer-Policy`），禁止访问隐藏文件。Docker 容器隔离，MySQL / Redis 仅绑定 `127.0.0.1`。HTTPS 由 Let's Encrypt + Certbot 自动管理。

**数据安全**：敏感配置通过环境变量注入，配置文件已加入 `.gitignore`，示例文件仅含占位符。

### 缓存架构

采用 Redis 多级缓存（Cache-Aside 模式），减少数据库查询：

| 缓存目标 | TTL | 说明 |
|----------|-----|------|
| 站点配置 | 24h | 每次页面加载都读，几乎不变 |
| 作者信息 | 24h | 单行数据，极少改动 |
| 分类列表 | 1h | 带聚合的 JOIN 查询，导航栏常驻 |
| 友链列表 | 24h | 极少变动 |
| 管理员权限 | 1h | 消除每次管理员请求的冗余查询 |
| 文章浏览量 | 实时 | Redis INCR 缓冲，定时批量同步 MySQL |

### 其他

- 全局健康检查接口（`/health`）
- 关键操作使用标准库控制台日志
- 文件上传与静态资源托管
- Docker Compose 一键部署（CI 构建镜像，服务器远程拉取）

## 前端视觉设计

自然系极简风格，薄荷绿（`#7ED7C1`）配色，暖白背景

## 技术栈

### 后端

| 组件 | 技术 |
|------|------|
| 语言 | [Go](https://golang.org/) 1.25 |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| ORM | [GORM](https://gorm.io/) v1.31 |
| 数据库 | MySQL 8.0+ |
| 缓存 | Redis（多级缓存 + 验证码 + 防刷限流 + 浏览量缓冲） |
| 日志 | 标准库 `log`（仅关键信息输出到控制台） |

### 前端

| 组件 | 技术 |
|------|------|
| 框架 | [Vue 3](https://vuejs.org/) + TypeScript |
| UI 组件库 | [Element Plus](https://element-plus.org/) |
| 状态管理 | [Pinia](https://pinia.vuejs.org/) |
| 路由 | [Vue Router](https://router.vuejs.org/) |
| HTTP 客户端 | [Axios](https://axios-http.com/) |
| 工具库 | [VueUse](https://vueuse.org/) |
| 构建工具 | [Vite](https://vitejs.dev/) |

### 部署

| 组件 | 技术 |
|------|------|
| 平台 | Linux |
| 反向代理 | Nginx |
| 容器化 | Docker Compose |
| CI/CD | GitHub Actions + 远程容器仓库 |

### Docker 挂载目录

项目通过 Docker Compose 部署时，以下目录会挂载到宿主机以实现数据持久化：

| 宿主机路径 | 容器内路径 | 用途 | 类型 |
|-----------|-----------|------|------|
| `backend/configs` | `/app/config` | 后端配置文件（yaml） | 绑定挂载 |
| `backend/uploads` | `/app/uploads` | 用户上传文件（头像、文章封面等） | 绑定挂载 |
| `frontend/public/author` | `/usr/share/nginx/html/author` | 作者头像与背景图 | 绑定挂载 |
| `logs` (Docker Volume) | `/app/log` | 后端日志 | 命名卷 |
| `mysql_data` (Docker Volume) | `/var/lib/mysql` | MySQL 数据 | 命名卷 |
| `redis_data` (Docker Volume) | `/data` | Redis 数据 | 命名卷 |

> 绑定挂载的目录直接映射到宿主机项目路径下，方便直接编辑配置或备份上传文件。命名卷由 Docker 自动管理。

## 快速开始

详见 [DEPLOY.md](DEPLOY.md) 与 [ARCHITECTURE.md](ARCHITECTURE.md)，涵盖本地开发与生产部署。

### 构建方式

推送到 `main` 分支时，GitHub Actions 会自动构建前后端镜像并推送到阿里云 ACR 与 Docker Hub，服务器直接拉取运行：

```bash
docker compose pull
docker compose up -d
```

## 开源协议

MIT License
