# 开发进度

## 后端

### 基础设施
- [x] 项目结构搭建
- [x] 配置加载（YAML + .env 环境切换）
- [x] 数据库连接池初始化（database/sql + go-sql-driver/mysql）
- [x] SQL 初始化脚本（user、category、article、comment）
- [x] Zap 日志初始化
- [x] 统一错误模型（AppError）
- [x] 统一 API 响应格式（ApiSuccessResponse / ApiErrorResponse）
- [x] JWT 服务（双 token：access + refresh）
- [x] JWT 鉴权中间件（RequireAuth / OptionalAuth）
- [x] 公共 RequireAdmin helper（common 包）
- [x] Context 常量统一管理（common 包）
- [x] user.Service 统一创建与注入（router.go）
- [x] CORS 中间件
- [x] Redis 连接池初始化（pkg/rdb）
- [x] Redis 验证码存取 & 防刷封装（pkg/rdb/verification）
- [x] SMTP 邮件发送封装（pkg/mail）
- [x] 验证码生成工具（pkg/util）
- [x] ErrTooManyRequests 错误码（HTTP 429）

### 用户模块（user）
- [x] Repository：GetUserByUsername / GetUserByEmail / GetUserByID / GetUserList / CreateUser / UpdateProfile / UpdatePassword / DeleteUser
- [x] Service：CreateUser / Login / GetByID / GetList / UpdateProfile / UpdatePassword / DeleteUser / CheckIsAdmin
- [x] Controller：Login / Register / SendCode / GetProfile / GetList / UpdateProfile / UpdatePassword / DeleteUser / RefreshToken
- [x] 路由注册与鉴权挂载
- [x] 邮箱验证码注册流程（发送验证码 → 校验验证码 → 创建用户）
- [x] 管理员权限校验（GetList / DeleteUser）
- [x] UserVO 响应过滤（隐藏 password）
- [x] initSystem 初始化管理员账号（用户名 admin，密码通过命令行参数传入）

### 文章模块（article）
- [x] Repository：GetArticleByID / GetArticleList / GetArticleByUserID / GetArticleByTitle / CreateArticle / UpdateArticle / DeleteArticle
- [x] Service：GetArticleList / GetArticleByID / GetArticleByUserID / GetArticleByTitle / CreateArticle / UpdateArticle / DeleteArticle
- [x] Controller：GetList / GetByID / CreateArticle / UpdateArticle / DeleteArticle
- [x] 路由注册（公开列表 + 详情，管理员增删改）

### 分类模块（category）
- [x] Repository：GetCategoryByID / GetCategoryByName / GetCategoryList / CreateCategory / UpdateCategory / DeleteCategory
- [x] Service：GetCategoryByID / GetCategoryByName / GetCategoryList / CreateCategory / UpdateCategory / DeleteCategory
- [x] Controller：GetList / Create / Update / Delete
- [x] 路由注册（公开列表，管理员增删改）

### 评论模块（comment）
- [x] Repository：GetCommentByID / GetCommentListByArticleID / CreateComment / DeleteComment
- [x] Service：GetCommentByID / GetCommentList / CreateComment / DeleteComment
- [x] Controller：GetList / CreateComment / DeleteComment（本人或管理员）
- [x] 路由注册（公开列表，登录用户创建/删除）

### 文件上传（uploads）
- [x] 上传接口
- [x] 静态文件服务

## 前端

### 基础设施
- [x] 项目初始化（Vue3 + TypeScript + Vite）
- [x] Element Plus UI 组件库
- [x] Axios 请求封装（src/api/request.ts）
- [x] Pinia 状态管理
- [x] Vue Router 路由配置
- [x] 页面布局组件（src/layouts）
- [x] API 接口封装（auth / user / article / comment）

### 页面
- [ ] 登录页面
- [ ] 注册页面（邮箱验证码）
- [ ] 首页
- [ ] 文章列表 / 详情
- [ ] 关于页面
- [ ] 后台管理面板
