# 第一阶段：构建 Go 应用
FROM golang:1.22 AS builder

# 设置国内 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . .

# 下载依赖并构建（关闭 CGO）
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o wechat ./main.go

# 第二阶段：运行时镜像
FROM alpine:3.19

WORKDIR /app

# 加速源 + 安装 ca-certificates + tzdata
RUN sed -i 's|http://dl-cdn.alpinelinux.org|https://mirrors.aliyun.com|g' /etc/apk/repositories && \
    apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

# 设置默认时区
ENV TZ=Asia/Shanghai

# 拷贝执行文件
COPY --from=builder /app/wechat .
COPY config/app.yml /app/config/app.yml

EXPOSE 8081
CMD ["/bin/sh", "-c", "/app/wechat migrate up && exec /app/wechat"]