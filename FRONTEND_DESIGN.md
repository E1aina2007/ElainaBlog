# ElainaBlog 前端设计文档

> 基于后端已完成的 API，描述前端的整体架构、页面、组件、状态管理与开发计划。

## 一、技术栈

| 类别 | 选型 |
|------|------|
| 框架 | Vue 3 + TypeScript |
| UI 组件库 | Element Plus |
| 状态管理 | Pinia |
| 路由 | Vue Router |
| HTTP 客户端 | Axios |
| 工具库 | VueUse |
| 构建工具 | Vite |
| Markdown 渲染 | markdown-it + highlight.js |
| Markdown 编辑器 | md-editor-v3（文章编写） |

---

## 二、目录结构

```
frontend/src/
├── api/                    # 后端接口封装
│   ├── request.ts          #   axios 实例、拦截器
│   ├── auth.ts             #   登录、注册、验证码、刷新 token
│   ├── user.ts             #   用户信息、修改资料、修改密码
│   ├── article.ts          #   文章 CRUD
│   └── comment.ts          #   评论 CRUD
├── components/             # 公共组件
│   ├── Navbar.vue          #   顶部导航栏
│   ├── Footer.vue          #   页脚
│   ├── ArticleCard.vue     #   文章卡片（列表页使用）
│   ├── Pagination.vue      #   分页组件
│   ├── MarkdownRenderer.vue#   Markdown 内容渲染
│   ├── CommentForm.vue     #   评论表单
│   └── CommentList.vue     #   评论列表
├── layouts/                # 布局
│   ├── DefaultLayout.vue   #   公开页面（导航栏 + 内容 + 页脚）
│   └── AdminLayout.vue     #   管理后台（侧栏 + 顶栏 + 内容）
├── router/
│   └── index.ts            # 路由配置 + 导航守卫
├── stores/                 # Pinia 状态管理
│   ├── user.ts             #   用户登录态、token、用户信息
│   └── app.ts              #   全局 UI 状态（侧栏折叠等）
├── utils/                  # 工具函数
│   ├── auth.ts             #   token 存取（localStorage）
│   ├── format.ts           #   日期格式化、文本截断
│   └── validate.ts         #   输入校验（正则）—— 已实现
├── styles/
│   └── global.css          # 全局样式、CSS 变量
├── views/                  # 页面视图
│   ├── Home.vue            #   首页（文章列表）
│   ├── ArticleDetail.vue   #   文章详情
│   ├── About.vue           #   关于页
│   ├── Login.vue           #   登录
│   ├── Register.vue        #   注册
│   └── admin/
│       ├── Dashboard.vue   #   管理仪表盘
│       ├── ArticleList.vue #   文章管理列表
│       ├── ArticleEdit.vue #   文章编辑（新建/编辑复用）
│       └── Profile.vue     #   个人资料 + 修改密码
├── App.vue                 # 根组件
└── main.ts                 # 入口（注册 Pinia、Router、Element Plus）
```

---

## 三、路由设计

| 路径 | 页面 | 布局 | 鉴权 |
|------|------|------|------|
| `/` | Home.vue | DefaultLayout | - |
| `/article/:id` | ArticleDetail.vue | DefaultLayout | - |
| `/about` | About.vue | DefaultLayout | - |
| `/login` | Login.vue | 无布局（全屏居中） | - |
| `/register` | Register.vue | 无布局（全屏居中） | - |
| `/admin` | Dashboard.vue | AdminLayout | ✅ 管理员 |
| `/admin/articles` | ArticleList.vue | AdminLayout | ✅ 管理员 |
| `/admin/article/create` | ArticleEdit.vue | AdminLayout | ✅ 管理员 |
| `/admin/article/edit/:id` | ArticleEdit.vue | AdminLayout | ✅ 管理员 |
| `/admin/profile` | Profile.vue | AdminLayout | ✅ 登录用户 |

### 路由守卫

```
router.beforeEach(to):
  1. 从 userStore 读取 token
  2. if to.meta.requiresAuth && 无 token → 重定向 /login
  3. if to.meta.requiresAdmin && !userStore.isAdmin → 重定向 /
  4. if (to 是 /login 或 /register) && 已登录 → 重定向 /
```

---

## 四、布局设计

### DefaultLayout — 公开页面

```
┌──────────────────────────────────────────┐
│  Navbar (logo / 首页 / 关于 / 登录按钮)    │
├──────────────────────────────────────────┤
│                                          │
│              <router-view />             │
│                                          │
├──────────────────────────────────────────┤
│  Footer (备案号 / 版权 / 社交链接)         │
└──────────────────────────────────────────┘
```

### AdminLayout — 管理后台

```
┌──────┬───────────────────────────────────┐
│      │  顶栏 (用户头像 / 退出登录)         │
│ 侧   ├───────────────────────────────────┤
│ 栏   │                                   │
│      │          <router-view />          │
│ 菜   │                                   │
│ 单   │                                   │
└──────┴───────────────────────────────────┘

侧栏菜单:
  - 仪表盘       /admin
  - 文章管理     /admin/articles
  - 个人资料     /admin/profile
```

### Login / Register — 全屏居中卡片，不套任何布局

---

## 五、状态管理

### stores/user.ts

```
state:
  accessToken: string            ← 与 localStorage 双向同步
  refreshToken: string           ← 与 localStorage 双向同步
  userInfo: {                    ← 登录后从 GET /user/profile 获取
    id: number
    username: string
    email: string
    avatar: string
    isAdmin: boolean
  } | null

getters:
  isLoggedIn → accessToken 非空
  isAdmin    → userInfo?.isAdmin === true

actions:
  login(email, password)       → POST /login → 存 token → fetchProfile()
  register(...)                → POST /register → 成功后跳转 /login
  logout()                     → 清空 token + userInfo → 跳转 /
  fetchProfile()               → GET /user/profile → 更新 userInfo
  refreshAccessToken()         → POST /refresh → 更新 accessToken
```

### stores/app.ts

```
state:
  sidebarCollapsed: boolean     ← 管理后台侧栏折叠
```

---

## 六、API 层设计

### api/request.ts — Axios 实例

```
baseURL: import.meta.env.VITE_API_BASE_URL  (默认 '/api/ui')

请求拦截器:
  if accessToken 存在 → Authorization: Bearer <token>

响应拦截器:
  成功 → 返回 response.data.data
  401 且有 refreshToken:
    → 调 refreshAccessToken()
    → 成功: 重发原请求
    → 失败: logout() + 跳 /login
  其他错误 → Promise.reject(格式化错误信息)
```

### 接口清单

#### api/auth.ts

| 函数 | 方法 | 路径 | 参数 |
|------|------|------|------|
| `login` | POST | `/login` | `{ email, password }` |
| `register` | POST | `/register` | `{ username, email, password, code }` |
| `sendCode` | POST | `/send-code` | `{ email }` |
| `refreshToken` | POST | `/refresh` | `{ refresh_token }` |

#### api/user.ts

| 函数 | 方法 | 路径 | 参数 |
|------|------|------|------|
| `getProfile` | GET | `/user/profile` | - |
| `updateProfile` | POST | `/user/profile` | `{ username, email, avatar }` |
| `updatePassword` | POST | `/user/password` | `{ old_password, new_password }` |
| `getUserList` | GET | `/user/list` | - |
| `deleteUser` | POST | `/user/delete` | `{ user_id }` |

#### api/article.ts

| 函数 | 方法 | 路径 | 参数 |
|------|------|------|------|
| `getArticleList` | GET | `/article/list` | `?page&pageSize&categoryId` |
| `getArticleDetail` | GET | `/article/:id` | - |
| `createArticle` | POST | `/article/create` | `{ title, content, ... }` |
| `updateArticle` | POST | `/article/update` | `{ id, title, content, ... }` |
| `deleteArticle` | POST | `/article/delete` | `{ id }` |

#### api/comment.ts

| 函数 | 方法 | 路径 | 参数 |
|------|------|------|------|
| `getComments` | GET | `/comment/:article_id` | - |
| `createComment` | POST | `/comment/create` | `{ article_id, content }` |
| `deleteComment` | POST | `/comment/delete` | `{ id }` |

---

## 七、页面详细设计

### 7.1 Login.vue — 登录

```
┌─────────────────────────┐
│       博客 Logo          │
│                         │
│   邮箱    [__________]  │
│   密码    [__________]  │
│                         │
│       [ 登  录 ]        │
│                         │
│   没有账号？去注册 →     │
└─────────────────────────┘
```

- 邮箱 `@blur` → `validateEmail()`，显示行内错误
- 密码仅判空，**不校验格式**（兼容旧密码）
- 提交 → `userStore.login()` → 成功跳 `/` 或来源页
- 后端 401 → 显示"邮箱或密码错误"

### 7.2 Register.vue — 注册

```
┌──────────────────────────────┐
│          博客 Logo            │
│                              │
│  用户名  [__________]        │
│          ⚠ 错误提示           │
│  邮箱    [__________] [发送]  │
│          ⚠ 错误提示           │
│  验证码  [______]             │
│  密码    [__________]        │
│          ⚠ 错误提示           │
│                              │
│        [ 注  册 ]            │
│                              │
│    已有账号？去登录 →         │
└──────────────────────────────┘
```

- 每个字段 `@blur` → 对应 `validate*()` → 显示/隐藏红色提示
- **[发送验证码]** 按钮：先 `validateEmail()` 通过 → 调 `sendCode()` → 60 秒倒计时禁用
- 提交前统一校验全部字段，任一不通过则阻止
- 成功 → `ElMessage.success('注册成功')` → 跳 `/login`

### 7.3 Home.vue — 首页

```
┌──────────────────────────────────────────┐
│  分类筛选 (el-tabs 或下拉)                │
├──────────────────────────────────────────┤
│  ArticleCard (置顶文章排在最前)            │
│  ArticleCard                             │
│  ArticleCard                             │
│  ...                                     │
├──────────────────────────────────────────┤
│            < 1  2  3  ... >              │
└──────────────────────────────────────────┘
```

- `onMounted` → `getArticleList({ page: 1, pageSize: 10 })`
- 分类切换 / 翻页 → 重新请求
- 置顶文章标记角标

### 7.4 ArticleDetail.vue — 文章详情

```
┌──────────────────────────────────────────┐
│  标题                                     │
│  作者 · 日期 · 浏览量 · 分类              │
├──────────────────────────────────────────┤
│  MarkdownRenderer 渲染 content            │
├──────────────────────────────────────────┤
│  评论区                                   │
│  CommentForm (未登录时提示"请先登录")       │
│  CommentList                             │
│    - 用户名 · 时间                        │
│      评论内容              [删除]         │
└──────────────────────────────────────────┘
```

- `onMounted` → `getArticleDetail(id)` + `getComments(id)`
- 评论删除按钮仅对本人或管理员可见

### 7.5 About.vue — 关于页

静态个人介绍页。可从 `site_config` 读取（昵称、工作、地址、社交链接），也可硬编码。

### 7.6 admin/Dashboard.vue — 仪表盘

```
┌────────┬────────┬────────┐
│ 文章数  │ 评论数  │ 用户数  │
└────────┴────────┴────────┘
```

统计卡片，MVP 阶段可先放占位，后续加统计 API。

### 7.7 admin/ArticleList.vue — 文章管理

```
[+ 新建文章]
┌──────────────────────────────────────────┐
│ 标题  │ 分类 │ 状态 │ 浏览 │ 日期 │ 操作  │  el-table
│ ...   │ ...  │ 草稿 │ 123  │ 04-29│ 编辑 删除 │
└──────────────────────────────────────────┘
分页
```

- 删除 → `ElMessageBox.confirm` 二次确认 → `deleteArticle()`

### 7.8 admin/ArticleEdit.vue — 文章编辑

```
标题    [____________________]
分类    [下拉选择]
摘要    [____________________]
封面    [上传图片]
内容    ┌─────────────────────┐
        │  Markdown 编辑器     │
        │  (md-editor-v3)     │
        └─────────────────────┘
        [存为草稿]  [发布]
```

- 路由 `/admin/article/create` → 新建模式
- 路由 `/admin/article/edit/:id` → 编辑模式，`onMounted` 拉取数据填充表单

### 7.9 admin/Profile.vue — 个人资料

```
头像      [上传]
用户名    [__________]    ← validateUsername
邮箱      [__________]    ← validateEmail
          [保存资料]

──────────────────
修改密码
旧密码    [__________]
新密码    [__________]    ← validatePassword
          [修改密码]
```

- 两个独立表单，分别提交
- `onMounted` → `getProfile()` 填充当前值

---

## 八、公共组件

| 组件 | 职责 |
|------|------|
| `Navbar.vue` | logo、导航链接(首页/关于)、右侧登录按钮或用户头像下拉菜单(个人中心/管理后台/退出) |
| `Footer.vue` | 备案号、版权年份、GitHub/Bilibili 等社交图标 |
| `ArticleCard.vue` | 水平卡片：左侧封面图、右侧标题+摘要+日期+分类标签，点击跳转详情 |
| `Pagination.vue` | 封装 `el-pagination`，props: `total` / `page` / `pageSize`，emit: `page-change` |
| `MarkdownRenderer.vue` | 接收 `content` prop，用 markdown-it 渲染为 HTML，highlight.js 高亮代码块 |
| `CommentForm.vue` | textarea + 提交按钮，未登录时显示提示文字并禁用 |
| `CommentList.vue` | 评论列表，每条显示头像、用户名、时间、内容，本人或管理员可见删除按钮 |

---

## 九、认证流程

```
用户登录
  │
  ▼
POST /login → 返回 accessToken + refreshToken
  │
  ▼
存入 localStorage + Pinia userStore
  │
  ▼
GET /user/profile → 存入 userStore.userInfo (含 isAdmin)
  │
  ▼
后续请求 → axios 拦截器自动注入 Authorization: Bearer <accessToken>
  │
  ▼
accessToken 过期 → 响应拦截器捕获 401
  │
  ▼
POST /refresh (带 refreshToken) → 获取新 accessToken → 重发原请求
  │
  ▼
refreshToken 也过期 → 清空状态 → 跳转 /login
```

---

## 十、输入校验规则

已在 `frontend/src/utils/validate.ts` 和 `backend/internal/user/validate.go` 中实现，前后端保持一致。

| 字段 | 规则 | 正则 |
|------|------|------|
| 邮箱 | 标准格式 local@domain.tld，≤100 字符 | `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$` |
| 用户名 | 中文/英文/数字/下划线，2-20 字符 | `^[\u4e00-\u9fa5a-zA-Z0-9_]{2,20}$` |
| 密码 | 8-72 字符，至少含一个字母和一个数字，仅允许合法字符 | 拆分为 4 条正则分别校验 |

### 各页面校验策略

| 页面 | 校验字段 | 时机 |
|------|----------|------|
| 注册 | 邮箱 + 用户名 + 密码 | @blur 实时 + 提交前统一 |
| 登录 | 邮箱（密码仅判空） | @blur + 提交前 |
| 发送验证码 | 邮箱 | 点击前校验 |
| 修改资料 | 用户名 + 邮箱 | @blur + 提交前 |
| 修改密码 | 新密码（旧密码仅判空） | @blur + 提交前 |

---

## 十一、样式方案

- **Element Plus** 作为核心 UI 组件库（表单、表格、消息提示、弹窗、菜单）
- `global.css` 定义 CSS 变量（主色、背景色、字体大小）
- 各 `.vue` 组件使用 `<style scoped>` 隔离样式
- 删除脚手架默认样式（`main.css` / `base.css` / `HelloWorld.vue` / `TheWelcome.vue` 等）
- 响应式：移动端 Navbar 折叠为汉堡菜单，文章列表单列，Admin 侧栏可折叠

---

## 十二、环境变量

`.env` 文件：

```
VITE_API_BASE_URL=/api/ui
```

开发时通过 `vite.config.ts` 的 `server.proxy` 代理到后端：

```ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

---

## 十三、建议开发顺序

| 阶段 | 内容 | 预计产出 |
|------|------|----------|
| **1. 基础骨架** | request.ts、stores/user.ts、router（含守卫）、两个 Layout、App.vue 改造、清理脚手架文件 | 可跑的空壳 |
| **2. 认证流程** | Login.vue、Register.vue、auth.ts API、token 刷新机制、Navbar 登录态切换 | 可注册 + 登录 + 自动刷新 |
| **3. 首页 + 详情** | Home.vue、ArticleCard、Pagination、ArticleDetail.vue、MarkdownRenderer、评论区 | 博客公开页面完整 |
| **4. 管理后台** | Dashboard、ArticleList、ArticleEdit（集成 Markdown 编辑器）、Profile | 后台可管理文章和个人信息 |
| **5. 收尾** | About 页、Footer、全局错误处理、响应式适配、loading 状态 | 上线就绪 |
