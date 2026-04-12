# ElainaWeb 基于 Gin 和 Vue3 + TypeScript 的个人小窝(网站)
遇见风的小窝 - 边学边搭建，搭到哪学到哪

## 技术栈

### 后端
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) v1.12
- **ORM**: [GORM](https://gorm.io) v1.31
- **数据库**: MySQL (驱动 gorm.io/driver/mysql v1.6)
- **日志**: [Zap](https://github.com/uber-go/zap) v1.27
- **配置管理**: [yaml](https://gopkg.in/yaml.v3) (YAML 文件读取/写入)
- **认证**: JWT (访问令牌 + 刷新令牌)
- **语言**: Go 1.25

### 前端
- **框架**: [Vue 3](https://vuejs.org/) + TypeScript
- **状态管理**: [Pinia](https://pinia.vuejs.org/)
- **路由**: [Vue Router](https://router.vuejs.org/)
- **构建工具**: [Vite](https://vitejs.dev/)

### 部署 & 运维
- **反向代理**: Nginx
- **容器化**: Docker + Docker Compose
- **日志管理**: Zap

## 项目结构
```
ElainaWeb/
├── backend/                  # 后端项目目录
│   ├── cmd/                  # 应用入口 
│   ├── config/               # 配置定义与加载
│   ├── global/               # 全局变量 
│   ├── internal/             # 业务逻辑
│   │   ├── handler/          # 请求处理器
│   │   ├── middleware/       # 中间件
│   │   ├── model/            # 数据模型
│   │   ├── repository/       # 数据访问层
│   │   ├── router/           # 路由注册
│   │   └── service/          # 业务服务层
│   ├── pkg/                  # 公共工具包
│   │   ├── db/               # 数据库初始化
│   │   ├── response/         # 统一响应封装
│   │   ├── util/             # 通用工具
│   │   └── zaplogger/        # Zap 日志初始化
│   ├── config.example.yaml   # 配置文件示例
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/                 # 前端项目目录
│   ├── src/
│   │   ├── api/              # 接口请求
│   │   ├── components/       # 公共组件
│   │   ├── layouts/          # 布局组件
│   │   ├── router/           # 路由配置
│   │   ├── stores/           # Pinia 状态管理
│   │   ├── styles/           # 全局样式
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
├── docker-compose.yml        # Docker 
└── README.md
```

## 待实现的功能
- [ ] 前端页面
- [ ] 后端基础框架
- [ ] 登录注册
- [ ] 用户信息管理/个人页面
- [ ] 文章管理
- [ ] 评论管理
