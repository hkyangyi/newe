# 构建阶段
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o newe-app .

# 运行阶段
FROM alpine:latest

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建应用用户
RUN adduser -D -g '' neweuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/newe-app .
COPY --from=builder /app/assets ./assets

# 创建必要的目录
RUN mkdir -p /app/assets/config \
    /app/assets/runtime \
    /app/upload/images \
    /app/upload/files

# 设置文件权限
RUN chown -R neweuser:neweuser /app
RUN chmod -R 755 /app

# 切换到非root用户
USER neweuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动应用
ENTRYPOINT ["./newe-app"]