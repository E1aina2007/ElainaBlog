# ElainaBlog

基于 Gin + Vue3 + TypeScript 的个人博客系统。

> 遇见风的小窝 —— 边学边搭版

## 技术栈

### 后端
- **语言**: Go 1.25
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) v1.12
- **数据库**: MySQL
- **缓存**: Redis（验证码存储 & 防刷限制）
- **日志**: [Zap](https://github.com/uber-go/zap) v1.27
- **认证**: JWT（双 token：access + refresh）
- **密码**: bcrypt
- **邮件**: net/smtp（SMTP 验证码发送）

### 前端
- **框架**: [Vue 3](https://vuejs.org/) + TypeScript
- **UI 组件库**: [Element Plus](https://element-plus.org/)
- **状态管理**: [Pinia](https://pinia.vuejs.org/)
- **路由**: [Vue Router](https://router.vuejs.org/)
- **HTTP 客户端**: [Axios](https://axios-http.com/)
- **工具库**: [VueUse](https://vueuse.org/)
- **构建工具**: [Vite](https://vitejs.dev/)

### 部署
- **平台**: Linux
- **反向代理**: Nginx
- **容器化**: Docker + Docker Compose

## 项目结构

```
ElainaBlog/
├── backend/
│   ├── cmd/                        # 应用入口（main、initSystem、runServer）
│   ├── config/                     # 配置定义与加载
│   │   └── db/                     # 数据库连接池 & SQL 初始化脚本
│   ├── internal/
│   │   ├── common/                 # JWT 服务、公共 helper、模型（AppError、ApiResponse）
│   │   ├── middleware/             # JWT 鉴权中间件
│   │   ├── router/                 # 路由注册
│   │   ├── user/                   # 用户模块（Controller / Service / Repository）
│   │   ├── article/                # 文章模块（Controller / Service / Repository）
│   │   ├── category/               # 分类模块（Controller / Service / Repository）
│   │   ├── comment/                # 评论模块（Controller / Service / Repository）
│   │   ├── upload/                 # 文件上传
│   │   └── site/                   # 站点配置
│   └── pkg/
│       ├── mail/                   # SMTP 邮件发送
│       ├── rdb/                    # Redis 连接 & 验证码操作
│       ├── util/                   # 工具函数（时间解析、验证码生成）
│       └── zaplogger/              # Zap 日志初始化
├── frontend/
│   └── src/
│       ├── api/                    # 后端接口封装（axios）
│       ├── components/             # 公共组件
│       ├── layouts/                # 页面布局
│       ├── router/                 # 路由配置
│       ├── stores/                 # Pinia 状态管理
│       ├── utils/                  # 工具函数
│       └── views/                  # 页面视图
├── docker-compose.yml
└── README.md
```

## API 路由

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/health` | - | 健康检查 |
| POST | `/api/ui/login` | - | 用户登录 |
| POST | `/api/ui/register` | - | 用户注册（需验证码） |
| POST | `/api/ui/send-code` | - | 发送邮箱验证码 |
| POST | `/api/ui/refresh` | - | 刷新 access token |
| GET | `/api/ui/user/profile` | ✅ | 获取当前用户信息 |
| GET | `/api/ui/user/list` | ✅ 管理员 | 用户列表 |
| POST | `/api/ui/user/profile` | ✅ | 修改个人资料 |
| POST | `/api/ui/user/password` | ✅ | 修改密码 |
| POST | `/api/ui/user/delete` | ✅ 管理员 | 删除用户 |
| GET | `/api/ui/category/list` | - | 分类列表 |
| POST | `/api/ui/category/create` | ✅ 管理员 | 创建分类 |
| POST | `/api/ui/category/update` | ✅ 管理员 | 更新分类 |
| POST | `/api/ui/category/delete` | ✅ 管理员 | 删除分类 |
| GET | `/api/ui/article/list` | - | 文章列表 |
| GET | `/api/ui/article/:id` | - | 文章详情 |
| POST | `/api/ui/article/create` | ✅ 管理员 | 创建文章 |
| POST | `/api/ui/article/update` | ✅ 管理员 | 更新文章 |
| POST | `/api/ui/article/delete` | ✅ 管理员 | 删除文章 |
| GET | `/api/ui/comment/:article_id` | - | 评论列表 |
| POST | `/api/ui/comment/create` | ✅ | 创建评论 |
| POST | `/api/ui/comment/delete` | ✅ | 删除评论（本人或管理员） |

## 快速启动

### 1. 配置

复制 `backend/config/config.example.yaml` 为 `backend/config/config.dev.yaml`，**务必修改以下示例内容**：

- **auth** — `access_token_secret` / `refresh_token_secret` 替换为你自己的随机密钥
- **db** — `username` / `password` 替换为你的 MySQL 账号密码
- **smtp** — `from` 替换为你的发件邮箱，`verification` 替换为邮箱 SMTP 授权码（如 QQ 邮箱需在设置中开启 SMTP 并获取授权码）
- **redis** — `address` / `password` 根据你的 Redis 实例修改
- **server** — `sessions_key` 替换为你自己的随机字符串

> ⚠️ `config.dev.yaml` 已加入 `.gitignore`，不会被提交到仓库。请勿将真实密码写入 `config.example.yaml`。

确保 `backend/config/.env` 中设置了运行模式：

```
MODE=dev
```

### 2. 初始化数据库

使用 `backend/config/db/initsql/0001_init_sql.sql` 在 MySQL 中创建数据库和表：

```bash
mysql -u root -p < backend/config/db/initsql/0001_init_sql.sql
```

### 3. 初始化管理员

```bash
cd backend
go run ./cmd initSystem <password>
```

> 必须提供管理员密码作为命令行参数，默认用户名为 `admin`。

### 4. 启动服务

```bash
cd backend
go run ./cmd runServer
```

### 5. 前端

```bash
cd frontend
npm install
npm run dev
```

> 前端默认运行在 `http://localhost:5173`
