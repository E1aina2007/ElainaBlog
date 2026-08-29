# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /elainablog ./cmd

# 运行阶段
FROM alpine:3.21

RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /elainablog .

# 创建上传目录，设置非 root 用户
RUN adduser -D -u 1001 appuser && \
    mkdir -p uploads && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 9178

CMD ["./elainablog"]
