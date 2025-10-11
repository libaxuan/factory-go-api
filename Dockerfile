# 多阶段构建 Dockerfile

# 第一阶段：构建
FROM golang:1.22-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && go mod verify

# 复制源代码
COPY main_multimodel.go ./
COPY config/ ./config/
COPY transformers/ ./transformers/
COPY config.json ./
COPY docs.html ./

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o factory-proxy main_multimodel.go

# 第二阶段：运行（Anthropic 原生模式）
FROM alpine:latest AS anthropic

# 安装 ca-certificates 用于 HTTPS 请求
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

WORKDIR /app

# 从构建阶段复制二进制文件和配置文件
COPY --from=builder /build/factory-proxy .
COPY --from=builder /build/config.json .
COPY --from=builder /build/docs.html .

# 设置文件权限
RUN chown -R app:app /app

# 切换到非 root 用户
USER app

# 暴露端口
EXPOSE 8000

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

# 启动命令
ENTRYPOINT ["./factory-proxy"]

# 第三阶段：运行（OpenAI 兼容模式）
FROM alpine:latest AS openai

# 安装 ca-certificates 用于 HTTPS 请求
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

WORKDIR /app

# 从构建阶段复制二进制文件和配置文件
COPY --from=builder /build/factory-proxy .
COPY --from=builder /build/config.json .
COPY --from=builder /build/docs.html .

# 设置文件权限
RUN chown -R app:app /app

# 切换到非 root 用户
USER app

# 暴露端口
EXPOSE 8003

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8003/v1/health || exit 1

# 启动命令
ENTRYPOINT ["./factory-proxy"]