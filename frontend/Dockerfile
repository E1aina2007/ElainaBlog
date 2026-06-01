# 构建阶段
FROM node:22-alpine AS builder

WORKDIR /app

# 限制 Node.js 堆内存为 384MB（实际 RSS 约 1.5 倍，~576MB）
ENV NODE_OPTIONS=--max-old-space-size=384

# 先复制依赖文件，利用 Docker 缓存
COPY package.json package-lock.json ./
RUN npm config set registry https://registry.npmmirror.com && npm ci

# 复制源码并分步构建
# 拆分为独立 RUN 以降低峰值内存：
#   npm ci 峰值 ~500MB，进程退出后释放
#   vue-tsc 峰值 ~400MB（384MB 堆 + 运行时开销）
#   vite build 峰值 ~300MB
COPY . .
RUN npx vue-tsc --build
RUN npx vite build

# 运行阶段
FROM nginx:stable-alpine

# 复制自定义 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 从构建阶段复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
