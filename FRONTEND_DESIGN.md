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
## 二、设计风格 — 清新治愈系

### 2.1 设计理念

以「春日清晨」为灵感，营造温暖、宁静、舒适的阅读氛围。通过柔和的色调、圆润的造型和细腻的动效，让每一位访客感受到如沐春风的治愈体验。

### 2.2 色彩系统

| 色彩角色 | 色值 | 用途 |
|----------|------|------|
| **主色调** | `#7ED7C1` | 按钮、链接、高亮、图标 |
| **主色调浅** | `#A8E6CF` | 悬停状态、渐变起点 |
| **主色调淡** | `#DCEDC8` | 背景装饰、标签底色 |
| **辅助色** | `#FFB7B2` | 次要按钮、爱心/点赞 |
| **强调色** | `#FFDAC1` | 提示信息、暖色点缀 |
| **背景主色** | `#FAFDF9` | 页面主背景（米白偏绿） |
| **背景次色** | `#F0F7F4` | 卡片、模块背景 |
| **文字主色** | `#2C3E50` | 标题、正文 |
| **文字次色** | `#5D6D7E` | 次要信息、描述 |
| **文字辅助** | `#95A5A6` | 时间、元信息 |
| **边框色** | `#E8F0ED` | 分割线、边框 |

### 2.3 字体规范

| 层级 | 字体 | 大小 | 字重 | 行高 |
|------|------|------|------|------|
| 页面标题 | Inter / 思源黑体 | 2.5rem (40px) | 600 | 1.3 |
| 文章标题 | Inter / 思源黑体 | 1.5rem (24px) | 600 | 1.4 |
| 正文 | Inter / 思源黑体 | 1rem (16px) | 400 | 1.8 |
| 小字说明 | Inter / 思源黑体 | 0.875rem (14px) | 400 | 1.6 |
| 标签/时间 | Inter / 思源黑体 | 0.75rem (12px) | 400 | 1.5 |

### 2.4 圆角与阴影

- **卡片圆角**：`16px`（大卡片）、`12px`（小卡片）、`8px`（按钮/标签）
- **柔和阴影**：`0 4px 20px rgba(126, 215, 193, 0.15)`（主色调投影）
- **悬浮阴影**：`0 8px 30px rgba(126, 215, 193, 0.25)`
- **内阴影**：`inset 0 2px 4px rgba(0,0,0,0.02)`（输入框）

### 2.5 动效规范

| 场景 | 时长 | 缓动函数 | 效果 |
|------|------|----------|------|
| 页面进入 | 400ms | `cubic-bezier(0.4, 0, 0.2, 1)` | opacity 0→1, translateY 20px→0 |
| 悬浮悬停 | 200ms | `ease-out` | scale 1→1.02, 阴影加深 |
| 按钮点击 | 100ms | `ease-in-out` | scale 1→0.98 |
| 卡片加载 | 300ms | `cubic-bezier(0.34, 1.56, 0.64, 1)` | stagger 50ms 依次出现 |
| 路由切换 | 300ms | `ease` | fade + slide |

### 2.6 装饰元素

- **自然插画**：角落的叶子、云朵、星星点缀（SVG 装饰，opacity 0.1-0.3）
- **渐变背景**：`linear-gradient(135deg, #FAFDF9 0%, #F0F7F4 100%)`
- **玻璃拟态**：关键卡片使用 `backdrop-filter: blur(10px)` + 半透明背景

---
## 三、目录结构

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

## 四、路由设计

| 路径 | 页面 | 布局 | 鉴权 |
|------|------|------|------|
| `/` | Home.vue | DefaultLayout | - |
| `/article/:id` | ArticleDetail.vue | DefaultLayout | - |
| `/about` | About.vue | DefaultLayout | - |
| `/author` | Author.vue | DefaultLayout | - |
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

## 五、布局设计

### DefaultLayout — 公开页面

```
┌──────────────────────────────────────────┐
│  Navbar (logo / 首页 / 关于 / 登录按钮)    │
├──────────────────────────────────────────┤
│                                          │
│              <router-view />             │
│                                          │
├──────────────────────────────────────────┤
│  Footer (备案号 / 版权)         │
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

## 六、状态管理

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

## 七、API 层设计

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
| `getAuthorInfo` | GET | `/author/info` | - | 获取作者公开信息 |
| `getAuthorStats` | GET | `/author/stats` | - | 获取作者统计数据 |

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

## 八、页面详细设计

### 8.1 Login.vue — 登录

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

### 8.2 Register.vue — 注册

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

### 8.3 Home.vue — 首页

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

### 8.4 ArticleDetail.vue — 文章详情

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

### 8.5 About.vue — 关于页

静态个人介绍页。展示博客的创建初衷、发展历程、技术栈说明等。

布局：
- 顶部装饰插图（治愈系插画）
- 博客故事区域（图文混排）
- 技术栈标签云
- 底部感谢语

---

### 8.5 Author.vue — 作者页

展示博主的个人信息，设计风格温馨治愈。

```
┌──────────────────────────────────────────┐
│                                          │
│   ╭──────────────────────────╮            │
│   │                          │            │
│   │     [头像 - 圆形大图]     │            │
│   │                          │            │
│   │      博主昵称             │            │
│   │   ✨ 一句话签名 ✨        │            │
│   │                          │            │
│   ╰──────────────────────────╯            │
│                                          │
│   ┌──────────┬──────────┬──────────┐     │
│   │  文章数   │  评论数   │  建站天数 │     │
│   │   128   │   256   │  365天   │     │
│   └──────────┴──────────┴──────────┘     │
│                                          │
│   ╭──────────────────────────────────╮   │
│   │  📍 位置    中国 · 杭州             │   │
│   │  💼 职业    前端开发者              │   │
│   │  🎯 爱好    摄影 · 旅行 · 阅读      │   │
│   │  📧 邮箱    hello@example.com       │   │
│   ╰──────────────────────────────────╯   │
│                                          │
│   个人简介                                │
│   ─────────────────────────────────────   │
│   这里是一段关于博主的温馨介绍文字，        │
│   讲述创作的初衷和对生活的热爱...          │
│                                          │
│   社交链接                                │
│   🐙 GitHub  📝 博客园  📺 Bilibili      │
│                                          │
└──────────────────────────────────────────┘
```

**设计细节**：
- 头像：`160px` 圆形，带 `4px` 白色边框 + 主色调阴影
- 统计卡片：玻璃拟态效果，数字使用主色调
- 信息项：每项带 emoji 图标，柔和分隔线
- 社交链接：圆形图标按钮，悬停时轻微上浮 + 变色

**数据结构**：
```ts
interface AuthorInfo {
  nickname: string
  avatar: string
  signature: string
  location: string
  occupation: string
  hobbies: string[]
  email: string
  bio: string
  social: {
    github?: string
    blog?: string
    bilibili?: string
    weibo?: string
  }
  stats: {
    articleCount: number
    commentCount: number
    daysSinceCreated: number
  }
}
```

### 8.6 admin/Dashboard.vue — 仪表盘

```
┌────────┬────────┬────────┐
│ 文章数  │ 评论数  │ 用户数  │
└────────┴────────┴────────┘
```

统计卡片，MVP 阶段可先放占位，后续加统计 API。

### 8.7 admin/ArticleList.vue — 文章管理

```
[+ 新建文章]
┌──────────────────────────────────────────┐
│ 标题  │ 分类 │ 状态 │ 浏览 │ 日期 │ 操作  │  el-table
│ ...   │ ...  │ 草稿 │ 123  │ 04-29│ 编辑 删除 │
└──────────────────────────────────────────┘
分页
```

- 删除 → `ElMessageBox.confirm` 二次确认 → `deleteArticle()`

### 8.8 admin/ArticleEdit.vue — 文章编辑

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

### 8.9 admin/Profile.vue — 个人资料

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

## 九、公共组件

| 组件 | 职责 |
|------|------|
| `Navbar.vue` | logo、导航链接(首页/关于)、右侧登录按钮或用户头像下拉菜单(个人中心/管理后台/退出) |
| `Footer.vue` | 备案号、版权年份、GitHub/Bilibili 等社交图标 |
| `ArticleCard.vue` | 水平卡片：左侧封面图、右侧标题+摘要+日期+分类标签，点击跳转详情 |
| `Pagination.vue` | 封装 `el-pagination`，props: `total` / `page` / `pageSize`，emit: `page-change` |
| `MarkdownRenderer.vue` | 接收 `content` prop，用 markdown-it 渲染为 HTML，highlight.js 高亮代码块 |
| `CommentForm.vue` | textarea + 提交按钮，未登录时显示提示文字并禁用 |
| `CommentList.vue` | 评论列表，每条显示头像、用户名、时间、内容，本人或管理员可见删除按钮 |
| `AuthorCard.vue` | 作者信息卡片（用于首页侧边栏或独立页面） |

---

## 十、认证流程

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

## 十一、输入校验规则

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

## 十二、样式方案

### 12.1 技术选型
- **Element Plus** 作为核心 UI 组件库（表单、表格、消息提示、弹窗、菜单）
- 通过 `ConfigProvider` 覆盖 Element Plus 主题色为清新治愈系
- 各 `.vue` 组件使用 `<style scoped>` 隔离样式
- 删除脚手架默认样式（`main.css` / `base.css` / `HelloWorld.vue` / `TheWelcome.vue` 等）

### 12.2 全局 CSS 变量 (`global.css`)

```css
:root {
  /* 主色调 */
  --primary: #7ED7C1;
  --primary-light: #A8E6CF;
  --primary-lighter: #DCEDC8;
  --primary-dark: #5BC4AD;
  
  /* 辅助色 */
  --accent: #FFB7B2;
  --accent-warm: #FFDAC1;
  --highlight: #FF9AA2;
  
  /* 背景色 */
  --bg-primary: #FAFDF9;
  --bg-secondary: #F0F7F4;
  --bg-card: #FFFFFF;
  --bg-glass: rgba(255, 255, 255, 0.8);
  
  /* 文字色 */
  --text-primary: #2C3E50;
  --text-secondary: #5D6D7E;
  --text-muted: #95A5A6;
  
  /* 边框与分割线 */
  --border: #E8F0ED;
  --divider: rgba(126, 215, 193, 0.2);
  
  /* 阴影 */
  --shadow-soft: 0 4px 20px rgba(126, 215, 193, 0.15);
  --shadow-hover: 0 8px 30px rgba(126, 215, 193, 0.25);
  --shadow-card: 0 2px 12px rgba(0, 0, 0, 0.04);
  
  /* 圆角 */
  --radius-lg: 16px;
  --radius-md: 12px;
  --radius-sm: 8px;
  --radius-full: 9999px;
  
  /* 字体 */
  --font-family: 'Inter', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  
  /* 过渡 */
  --transition-fast: 150ms ease;
  --transition-base: 300ms ease;
  --transition-slow: 500ms ease;
}
```

### 12.3 Element Plus 主题覆盖

```ts
// main.ts
import { ElConfigProvider } from 'element-plus'

const themeConfig = {
  button: {
    color: '#7ED7C1',
    hoverColor: '#5BC4AD',
    activeColor: '#4BA898',
    borderRadius: '8px',
  },
  menu: {
    activeColor: '#7ED7C1',
    hoverBgColor: '#F0F7F4',
  },
  card: {
    borderRadius: '16px',
  },
  input: {
    borderRadius: '8px',
    focusBorderColor: '#7ED7C1',
  },
  tag: {
    borderRadius: '6px',
  },
}
```

### 12.4 响应式断点

| 断点 | 宽度 | 布局调整 |
|------|------|----------|
| `xs` | < 576px | 单列布局，Navbar 折叠为汉堡菜单 |
| `sm` | ≥ 576px | 保留基础边距 |
| `md` | ≥ 768px | 文章列表双列 |
| `lg` | ≥ 992px | 完整布局，侧边栏显示 |
| `xl` | ≥ 1200px | 最大内容宽度 `1200px` 居中 |

### 12.5 动画工具类

```css
/* 页面进入动画 */
.fade-up-enter {
  animation: fadeUp 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 卡片依次出现 */
.stagger-enter > *:nth-child(1) { animation-delay: 0ms; }
.stagger-enter > *:nth-child(2) { animation-delay: 50ms; }
.stagger-enter > *:nth-child(3) { animation-delay: 100ms; }
/* ... */

/* 悬浮效果 */
.hover-lift {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.hover-lift:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-hover);
}
```

---

## 十三、环境变量

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

## 十四、建议开发顺序

| 阶段 | 内容 | 预计产出 |
|------|------|----------|
| **1. 基础骨架** | request.ts、stores/user.ts、router（含守卫）、两个 Layout、App.vue 改造、清理脚手架文件 | 可跑的空壳 |
| **2. 认证流程** | Login.vue、Register.vue、auth.ts API、token 刷新机制、Navbar 登录态切换 | 可注册 + 登录 + 自动刷新 |
| **3. 首页 + 详情** | Home.vue、ArticleCard、Pagination、ArticleDetail.vue、MarkdownRenderer、评论区 | 博客公开页面完整 |
| **4. 管理后台** | Dashboard、ArticleList、ArticleEdit（集成 Markdown 编辑器）、Profile | 后台可管理文章和个人信息 |
| **5. 收尾** | About 页、Footer、全局错误处理、响应式适配、loading 状态 | 上线就绪 |
