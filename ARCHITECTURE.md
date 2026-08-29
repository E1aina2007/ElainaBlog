# 项目架构文档

## 项目结构

```
ElainaBlog/
├── backend/                    # Go 后端服务
│   ├── Dockerfile              # 多阶段构建（Go 编译 + Alpine 运行）
│   ├── cmd/                    # 应用入口与 HTTP 服务启动
│   ├── configs/                # 外部配置（yaml；环境变量文件为 backend/.env，均不提交 Git）
│   ├── db/
│   │   └── migrations/         # SQL 迁移脚本（.up.sql / .down.sql）
│   ├── internal/               # 核心业务逻辑（不可被外部导入）
│   │   ├── config/             # 配置定义与加载
│   │   ├── db/                 # 数据库连接池
│   │   ├── mail/               # SMTP 邮件发送与邮箱哈希
│   │   ├── util/               # 工具子包（timeutil / verifycode）
│   │   ├── auth/               # JWT 服务与 token 管理
│   │   ├── response/           # API 统一响应与业务错误
│   │   ├── middleware/         # Gin 中间件
│   │   │   ├── jwt/            # JWT 鉴权与管理员权限
│   │   │   ├── ratelimit/      # 接口速率限制
│   │   │   ├── redis/          # Redis 连接与操作
│   │   │   └── uploadlimit/    # 上传频率限制
│   │   ├── router/             # 路由注册
│   │   ├── user/               # 用户模块（Controller / Service / Repository）
│   │   ├── article/            # 文章模块
│   │   ├── category/           # 分类模块
│   │   ├── comment/            # 评论模块
│   │   ├── site/               # 站点配置模块
│   │   └── upload/             # 文件上传模块
├── frontend/                   # Vue3 前端应用
│   ├── Dockerfile              # 多阶段构建（Node 构建 + Nginx 托管）
│   ├── nginx.conf              # 容器内 Nginx 配置（SPA fallback + 反向代理）
│   ├── .dockerignore
│   ├── public/
│   │   └── author/             # 作者头像与背景图（不提交 Git）
│   └── src/
│       ├── api/                # 后端接口封装（axios）
│       ├── components/         # 公共组件
│       ├── layouts/            # 页面布局（DefaultLayout / AdminLayout）
│       ├── router/             # 路由配置
│       ├── stores/             # Pinia 状态管理
│       ├── styles/             # 全局样式、CSS 变量
│       ├── utils/              # 工具函数
│       └── views/              # 页面视图
├── docs/                       # 设计与规划文档（TODO / issues / plans）
├── docker-compose.yml          # Docker Compose 配置（frontend + backend + mysql + redis）
└── README.md                   # 项目说明
```

## API 路由

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/health` | - | 健康检查 |
| POST | `/api/ui/login` | - | 用户登录 |
| POST | `/api/ui/register` | - | 用户注册（需验证码） |
| POST | `/api/ui/send-code` | - | 发送邮箱验证码 |
| POST | `/api/ui/refresh` | - | 刷新 access token |
| GET | `/api/ui/user/profile` | 用户 | 获取当前用户信息 |
| GET | `/api/ui/user/list` | 管理员 | 用户列表 |
| POST | `/api/ui/user/profile` | 用户 | 修改个人资料 |
| POST | `/api/ui/user/password` | 用户 | 修改密码 |
| POST | `/api/ui/user/delete` | 管理员 | 删除用户 |
| GET | `/api/ui/category/list` | - | 分类列表 |
| POST | `/api/ui/category/create` | 管理员 | 创建分类 |
| POST | `/api/ui/category/update` | 管理员 | 更新分类 |
| POST | `/api/ui/category/delete` | 管理员 | 删除分类 |
| GET | `/api/ui/article/list` | - | 文章列表 |
| GET | `/api/ui/article/:id` | - | 文章详情 |
| POST | `/api/ui/article/create` | 管理员 | 创建文章 |
| POST | `/api/ui/article/update` | 管理员 | 更新文章 |
| POST | `/api/ui/article/delete` | 管理员 | 删除文章 |
| GET | `/api/ui/comment/:article_id` | - | 评论列表 |
| POST | `/api/ui/comment/create` | 用户 | 创建评论 |
| POST | `/api/ui/comment/delete` | 用户 | 删除评论 |
| GET | `/api/ui/site` | - | 站点配置信息 |
| POST | `/api/ui/site` | 用户 | 更新站点配置 |
| GET | `/api/ui/author/info` | - | 作者公开信息 |
| GET | `/api/ui/author/stats` | - | 作者统计数据 |
| GET | `/api/ui/dashboard/stats` | 管理员 | 仪表盘统计 |
| POST | `/api/ui/upload` | 用户 | 文件上传 |

## 输入校验规则

前后端统一校验规则：

| 字段 | 规则 |
|------|------|
| 邮箱 | 标准格式，≤100 字符 |
| 用户名 | 中文/英文/数字/下划线，2-20 字符 |
| 密码 | 8-72 字符，至少含一个字母和一个数字 |

## 安全设计

- **密码存储**: bcrypt 哈希，自动加盐
- **JWT 密钥**: Access/Refresh 使用不同密钥
- **SQL 注入**: 使用参数化查询
- **CORS**: 可配置允许来源
- **验证码**: Redis 存储，防重发机制（60秒间隔）

## Docker 部署架构

```
用户 → 宿主机 Nginx (:80/:443)
        ├── /api/, /uploads/, /health → backend 容器 (:9178)
        └── 其余请求 → frontend 容器 (:3000)
                         └── 容器内 Nginx 反向代理 API → backend:9178
```

| 容器 | 端口 | 说明 |
|------|------|------|
| `elainablog-frontend` | 127.0.0.1:3000:80 | Nginx 托管前端静态文件 |
| `elainablog-backend` | 127.0.0.1:9178:9178 | Go API 服务 |
| `elainablog-migrate` | 内部 | golang-migrate 数据库迁移（一次性执行） |
| `elainablog-mysql` | 内部 | MySQL 8.0 |
| `elainablog-redis` | 内部 | Redis 7 |

### 数据卷

| 类型 | 宿主机路径 | 容器路径 | 说明 |
|------|-----------|---------|------|
| 命名卷 | `mysql_data` | /var/lib/mysql | MySQL 数据持久化 |
| 命名卷 | `redis_data` | /data | Redis 数据持久化 |
| 绑定挂载 | `backend/uploads/` | /app/uploads | 上传文件持久化 |
| 绑定挂载 | `frontend/public/author/` | /usr/share/nginx/html/author | 作者头像与背景图 |
| 绑定挂载 | `backend/configs/` | /app/config | 后端配置文件（yaml） |

> SQL 迁移脚本挂载到 `migrate` 容器，由 golang-migrate 在 MySQL 就绪后自动执行；后端镜像不再内置迁移。
