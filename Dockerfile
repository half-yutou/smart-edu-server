# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 Go 代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译应用
# CGO_ENABLED=0 确保生成静态链接的二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装基础依赖（如时区数据）
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]