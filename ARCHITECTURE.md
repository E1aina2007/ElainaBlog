# 项目架构文档

## 项目结构

```
ElainaBlog/
├── backend/                    # Go 后端服务
│   ├── Dockerfile              # 多阶段构建（Go 编译 + Alpine 运行）
│   ├── cmd/                    # 应用入口（main、initSystem、runServer）
│   ├── config/                 # 配置定义与加载（Go 源码 + 本地开发配置）
│   │   └── db/                 # 数据库连接池 & SQL 迁移脚本
│   ├── internal/               # 核心业务逻辑（不可被外部导入）
│   │   ├── common/             # JWT 服务、公共 helper、模型
│   │   ├── middleware/         # JWT 鉴权中间件
│   │   ├── router/             # 路由注册
│   │   ├── user/               # 用户模块（Controller / Service / Repository）
│   │   ├── article/            # 文章模块
│   │   ├── category/           # 分类模块
│   │   ├── comment/            # 评论模块
│   │   ├── site/               # 站点配置模块
│   │   └── upload/             # 文件上传模块
│   └── pkg/                    # 公共工具包
│       ├── mail/               # SMTP 邮件发送
│       ├── rdb/                # Redis 连接 & 验证码操作
│       ├── util/               # 工具函数
│       └── zaplogger/          # Zap 日志初始化
├── config/                     # 外部配置目录（Docker 挂载，不提交 Git）
│   ├── .env.example            # Docker Compose 环境变量示例
│   └── backend/
│       ├── .env.example        # 运行模式示例（MODE=dev/prod）
│       └── config.example.yaml # 后端配置示例
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
├── docs/                       # 设计文档
│   ├── BACKEND_DESIGN.md       # 后端设计文档
│   ├── FRONTEND_DESIGN.md      # 前端设计文档
│   └── PROGRESS.md             # 开发进度
├── nginx.conf                  # 宿主机 Nginx 入口配置
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
| `elainablog-mysql` | 内部 | MySQL 8.0 |
| `elainablog-redis` | 内部 | Redis 7 |

### 数据卷

| 类型 | 宿主机路径 | 容器路径 | 说明 |
|------|-----------|---------|------|
| 命名卷 | `mysql_data` | /var/lib/mysql | MySQL 数据持久化 |
| 命名卷 | `redis_data` | /data | Redis 数据持久化 |
| 命名卷 | `uploads` | /app/uploads | 上传文件持久化 |
| 命名卷 | `logs` | /app/log | 后端日志持久化 |
| 绑定挂载 | `frontend/public/author/` | /usr/share/nginx/html/author | 作者头像与背景图 |
| 绑定挂载 | `config/backend/` | /app/config | 后端配置文件（.env + yaml） |

> SQL 迁移脚本烘焙在后端镜像的 `/app/migrations/` 路径，不被 config 挂载覆盖。
