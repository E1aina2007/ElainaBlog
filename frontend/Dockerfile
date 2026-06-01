# 构建阶段
FROM node:22-alpine AS builder

WORKDIR /app

# 限制 Node.js 堆内存，防止 OOM
ENV NODE_OPTIONS=--max-old-space-size=512

# 先复制依赖文件，利用 Docker 缓存
COPY package.json package-lock.json ./
RUN npm config set registry https://registry.npmmirror.com && npm ci

# 复制源码并构建
# 不使用 npm run build（会并行执行 type-check 和 build，内存翻倍）
# 改为串行执行，降低峰值内存占用
COPY . .
RUN npx vue-tsc --build && npx vite build

# 运行阶段
FROM nginx:stable-alpine

# 复制自定义 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 从构建阶段复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
