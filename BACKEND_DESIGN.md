# ElainaBlog 后端设计文档

## 1. 技术栈概览

| 层级 | 技术选型 | 说明 |
|------|----------|------|
| 语言 | Go 1.25+ | 强类型、高性能、并发友好 |
| Web框架 | Gin | 轻量级高性能HTTP框架 |
| 数据库 | MySQL | 关系型数据持久化 |
| 缓存 | Redis | 验证码存储、会话管理 |
| 日志 | Zap + Lumberjack | 高性能结构化日志，支持日志轮转 |
| 认证 | JWT (jwt/v5) | 双Token机制（Access + Refresh） |
| 邮件 | SMTP | 验证码邮件发送 |
| 配置 | YAML + godotenv | 环境分离的配置管理 |

## 2. 项目结构

采用 **Clean Architecture / 分层架构** 思想组织代码：

```
backend/
├── cmd/                    # 应用入口
│   ├── main.go            # 初始化流程、命令路由
│   └── handler.go         # 子命令实现（initSystem / runServer）
│
├── config/                # 配置管理
│   ├── config.go          # 全局配置结构体、加载逻辑
│   ├── DbConfig.go        # 数据库配置
│   ├── AuthConfig.go      # JWT认证配置
│   ├── RedisConfig.go     # Redis配置
│   ├── ServerConfig.go    # HTTP服务器配置
│   ├── ZapConfig.go       # 日志配置
│   ├── UploadConfig.go    # 文件上传配置
│   ├── smtpConfig.go      # SMTP配置
│   ├── verificationConfig.go # 验证码配置
│   └── db/
│       └── mysql.go       # 数据库连接初始化
│
├── internal/              # 核心业务逻辑（不可被外部导入）
│   ├── router/            # 路由注册
│   │   └── router.go      # 统一注册所有API路由
│   │
│   ├── middleware/        # HTTP中间件
│   │   └── jwt_auth.go    # JWT鉴权中间件
│   │
│   ├── common/            # 公共组件
│   │   ├── jwt_auth_service.go  # JWT服务实现
│   │   ├── context_keys.go      # Context键定义
│   │   ├── require_admin.go   # 管理员权限检查
│   │   └── model/
│   │       ├── api_response.go  # 统一响应结构
│   │       └── app_error.go     # 业务错误定义
│   │
│   ├── user/              # 用户模块
│   │   ├── controller.go    # HTTP处理层
│   │   ├── service.go       # 业务逻辑层
│   │   ├── repository.go    # 数据访问层
│   │   └── validate.go      # 参数校验
│   │
│   ├── article/           # 文章模块
│   │   ├── controller.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── category/          # 分类模块
│   │   ├── controller.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── comment/           # 评论模块
│   │   ├── controller.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── site/              # 站点配置模块
│   │   ├── controller.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   └── upload/            # 文件上传模块
│       ├── controller.go
│       └── storage.go
│
└── pkg/                   # 公共工具包（可被外部导入）
    ├── mail/              # 邮件发送
    ├── rdb/               # Redis客户端封装
    ├── util/              # 通用工具函数
    └── zaplogger/         # Zap日志初始化
```

## 3. 架构分层设计

### 3.1 分层职责

```
┌─────────────────────────────────────────┐
│  Controller (控制器层)                   │
│  - 处理HTTP请求/响应                     │
│  - 参数绑定与校验                        │
│  - 调用Service层                        │
├─────────────────────────────────────────┤
│  Middleware (中间件层)                  │
│  - JWT鉴权                              │
│  - CORS跨域                             │
│  - 请求日志                             │
├─────────────────────────────────────────┤
│  Service (业务逻辑层)                   │
│  - 业务规则实现                          │
│  - 事务协调                              │
│  - 权限校验                              │
├─────────────────────────────────────────┤
│  Repository (数据访问层)                   │
│  - SQL查询封装                           │
│  - 数据模型映射                          │
│  - 数据库事务                            │
├─────────────────────────────────────────┤
│  Model/Entity (数据模型层)               │
│  - 领域对象定义                          │
│  - DTO/VO转换                            │
└─────────────────────────────────────────┘
```

### 3.2 依赖方向

**严格遵循依赖倒置原则**：
- `Controller` → `Service` → `Repository`
- 上层通过接口/构造函数依赖下层，不允许跨层调用
- 各模块独立，通过 `router.go` 统一组装

## 4. 核心模块设计

### 4.1 认证授权体系

#### 双Token机制

```go
// TokenClaims JWT声明结构
type TokenClaims struct {
    UserID    int64  `json:"user_id"`
    TokenType string `json:"token_type"`  // "access" | "refresh"
    jwt.RegisteredClaims
}
```

| Token类型 | 有效期 | 用途 |
|-----------|--------|------|
| Access Token | 2小时 | 日常接口鉴权，频繁使用 |
| Refresh Token | 7天 | 换取新的Access Token，长期有效 |

**流程**：
1. 登录成功后返回双Token
2. 前端将Access Token放入请求头 `Authorization: Bearer <token>`
3. Access Token过期后，使用Refresh Token调用 `/api/ui/refresh` 换取新的Access Token
4. Refresh Token也过期则需重新登录

#### JWT中间件

```go
// RequireAuth 强制鉴权
// - 提取Authorization头
// - 解析并验证Access Token
// - 将UserID写入Context供后续使用
auth := middleware.NewJwtAuthMiddleware(common.JwtAuth)
r.Group("/api/ui/user", auth.RequireAuth())

// OptionalAuth 可选鉴权
// - 有Token则解析，无Token也放行
// - 用于部分公开接口识别用户身份
```

### 4.2 错误处理体系

#### 业务错误码规范

```
400xxx - 客户端参数错误 → HTTP 400
401xxx - 未认证          → HTTP 401
403xxx - 无权限          → HTTP 403
404xxx - 资源不存在      → HTTP 404
409xxx - 资源冲突        → HTTP 409
429xxx - 请求过于频繁    → HTTP 429
500xxx - 服务器内部错误  → HTTP 500
```

#### 统一响应格式

```go
type ApiResponse struct {
    Success   bool   `json:"success"`      // 是否成功
    Code      int    `json:"code"`         // 业务错误码
    Message   string `json:"message"`      // 提示信息
    Data      any    `json:"data"`         // 业务数据
    Error     any    `json:"error"`        // 错误详情
    Timestamp int64  `json:"timestamp"`    // 毫秒时间戳
}

// 成功示例
{"success":true,"code":0,"message":"success","data":{"user_id":1},"timestamp":1704067200000}

// 失败示例
{"success":false,"code":401001,"message":"unauthorized","error":"token已过期","timestamp":1704067200000}
```

### 4.3 用户模块

| 功能 | 接口 | 权限 |
|------|------|------|
| 注册 | POST /api/ui/register | 公开 |
| 登录 | POST /api/ui/login | 公开 |
| 刷新Token | POST /api/ui/refresh | 公开（需Refresh Token）|
| 发送验证码 | POST /api/ui/send-code | 公开 |
| 获取个人信息 | GET /api/ui/user/profile | 登录用户 |
| 更新个人信息 | POST /api/ui/user/profile | 登录用户 |
| 修改密码 | POST /api/ui/user/password | 登录用户 |
| 获取用户列表 | GET /api/ui/user/list | 管理员 |
| 删除用户 | POST /api/ui/user/delete | 管理员 |

**密码安全**：
- 使用 `bcrypt` 算法哈希存储（默认cost）
- 登录时 `bcrypt.CompareHashAndPassword` 比对

**注册验证**：
- 邮箱验证码通过Redis存储，有效期5分钟
- 验证码重发间隔60秒防刷

### 4.4 文章模块

| 功能 | 接口 | 权限 |
|------|------|------|
| 获取文章列表 | GET /api/ui/article/list | 公开 |
| 获取文章详情 | GET /api/ui/article/:id | 公开 |
| 创建文章 | POST /api/ui/article/create | 管理员 |
| 更新文章 | POST /api/ui/article/update | 管理员 |
| 删除文章 | POST /api/ui/article/delete | 管理员 |

**文章状态**：
- `is_top`: 是否置顶
- `is_draft`: 是否草稿（草稿不对外展示）

### 4.5 评论模块

| 功能 | 接口 | 权限 |
|------|------|------|
| 获取文章评论 | GET /api/ui/comment/:article_id | 公开 |
| 发表评论 | POST /api/ui/comment/create | 登录用户 |
| 删除评论 | POST /api/ui/comment/delete | 登录用户 |

### 4.6 文件上传

- 本地存储模式，存储路径可配置
- 上传文件大小限制可配置
- 通过 `/uploads/` 静态路径访问

## 5. 数据持久化

### 5.1 数据库连接池

```go
db.SetMaxIdleConns(dbConfig.MaxIdleConns)  // 最大空闲连接
db.SetMaxOpenConns(dbConfig.MaxOpenConns)  // 最大打开连接
```

### 5.2 DSN安全生成

使用 `mysql.Config` 结构体而非字符串拼接，自动处理特殊字符转义：

```go
cfg := mysql.Config{
    User:   m.Username,
    Passwd: m.Password,
    Net:    "tcp",
    Addr:   net.JoinHostPort(m.Host, strconv.Itoa(m.Port)),
    DBName: m.DBName,
    // ...
}
return cfg.FormatDSN()
```

## 6. 配置管理

### 6.1 环境分离

```
.env              # 环境变量（确定MODE）
config.dev.yaml   # 开发环境配置
config.prod.yaml  # 生产环境配置（未提供示例）
```

### 6.2 配置加载流程

1. 从 `.env` 读取 `MODE` 环境变量
2. 根据 `MODE` 选择对应YAML配置文件
3. 解析YAML到全局 `GlobalConfig` 变量
4. 设置默认值（如Token有效期、Issuer等）

## 7. 启动流程

```
init()
  ├── 加载配置文件 (config.LoadConfigFromYml)
  ├── 初始化Zap日志 (zaplogger.InitLogger)
  ├── 初始化数据库 (db.InitDB)
  ├── 初始化Redis (rdb.InitRedis)
  └── 初始化JWT服务 (common.InitJwtAuth)

main()
  └── 根据命令参数执行：
      ├── initSystem <password>  # 创建管理员账号
      └── runServer              # 启动HTTP服务

runServer()
  ├── 设置Gin运行模式
  ├── 注册中间件（Recovery、CORS）
  ├── 注册路由 (router.RouterInit)
  ├── 启动HTTP服务（goroutine）
  └── 监听系统信号，优雅关停
```

## 8. 优雅关停

```go
// 监听系统信号
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
sig := <-quit

// 10秒超时上下文
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 优雅关闭HTTP服务
server.Shutdown(ctx)
```

## 9. 安全设计

- **密码存储**：bcrypt哈希，自动加盐
- **JWT密钥**：Access/Refresh使用不同密钥，配置分离
- **SQL注入**：使用标准库 `database/sql` 参数化查询
- **CORS**：可配置允许来源，默认仅允许开发环境前端地址
- **验证码**：Redis存储，防重发机制

## 10. 扩展性考虑

- 模块化设计：每个业务模块（user/article/comment等）独立，新增模块只需在router注册
- 依赖注入：通过构造函数注入依赖，便于单元测试和Mock
- 配置驱动：各项参数均可通过YAML配置，无需改代码
- 错误码预留：按区间划分错误码，便于后续扩展
