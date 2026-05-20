# ElainaBlog

> 基于 Gin + Vue 3 二次开发的个人博客网站框架
>
> *Elaina 这个名字来源于动漫《魔女之旅》中的主角灰之魔女伊蕾娜。她很可爱。*

## 功能

### 文章系统

- Markdown 编辑与渲染，支持代码高亮（highlight.js）、表格、任务列表等扩展语法
- 文章分类管理，首页按分类筛选
- 文章置顶
- 草稿与发布状态切换
- 封面图上传
- 阅读量统计

### 用户系统

- 邮箱注册 + 邮箱验证码验证
- JWT 双 Token 认证（Access Token + Refresh Token），无感刷新
- 个人资料编辑与头像上传
- 密码修改

### 评论系统

- 文章评论与回复
- 登录用户可发表、删除自己的评论
- 管理员可管理全部评论

### 留言板

- 公开留言列表
- 登录用户可发表、删除留言

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

### 安全与限流

- 基于 IP 的请求频率限制（防刷）
- 验证码发送间隔限制
- IP 封禁机制
- 密码 bcrypt 加密存储
- 角色权限控制（普通用户 / 管理员）

### 其他

- 全局健康检查接口（`/health`）
- 日志按级别分文件输出（Zap + Lumberjack 自动轮转）
- 文件上传与静态资源托管
- Docker Compose 一键部署

## 前端视觉设计

自然系极简风格，薄荷绿（`#7ED7C1`）+ 柔粉（`#FFB7B2`）配色，暖白背景。毛玻璃导航栏与统计面板，花瓣飘落动画，文章卡片交错淡入上浮，768px 响应式断点，Inter + PingFang SC 字体。

## 技术栈

### 后端

| 组件 | 技术 |
|------|------|
| 语言 | [Go](https://golang.org/) 1.25 |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| 数据库 | MySQL 8.0+ |
| 缓存 | Redis（验证码存储 & 防刷限制） |
| 日志 | [Zap](https://github.com/uber-go/zap) v1.27 + Lumberjack |

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

### Docker 挂载目录

项目通过 Docker Compose 部署时，以下目录会挂载到宿主机以实现数据持久化：

| 宿主机路径 | 容器内路径 | 用途 | 类型 |
|-----------|-----------|------|------|
| `config/backend` | `/app/config` | 后端配置文件（`.env` 和 `yaml`） | 绑定挂载 |
| `backend/uploads` | `/app/uploads` | 用户上传文件（头像、文章封面等） | 绑定挂载 |
| `frontend/public/author` | `/usr/share/nginx/html/author` | 作者头像与背景图 | 绑定挂载 |
| `logs` (Docker Volume) | `/app/log` | 后端日志 | 命名卷 |
| `mysql_data` (Docker Volume) | `/var/lib/mysql` | MySQL 数据 | 命名卷 |
| `redis_data` (Docker Volume) | `/data` | Redis 数据 | 命名卷 |

> 绑定挂载的目录直接映射到宿主机项目路径下，方便直接编辑配置或备份上传文件。命名卷由 Docker 自动管理。

## 快速开始

详见 【快速开始】[DEPLOY.md](DEPLOY.md)与【项目架构】[ARCHITECTURE.md](ARCHITECTURE.md)，涵盖本地开发与生产部署。

## 开源协议

MIT License
