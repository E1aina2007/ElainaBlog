# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并构建
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /elainablog ./cmd

# 运行阶段
FROM alpine:3.21

RUN apk add --no-cache tzdata mysql-client
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /elainablog .
# 复制数据库迁移脚本（独立路径，不被 config 挂载覆盖）
COPY --from=builder /app/config/db/SQLscript/ /app/migrations/

# 创建上传目录和日志目录，设置非 root 用户
RUN adduser -D -u 1001 appuser && \
    mkdir -p uploads log && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 9178

CMD ["./elainablog", "runServer"]
