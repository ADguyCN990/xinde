# --- Stage 1: Build ---
# 使用官方的 Go 镜像作为构建环境
# 选择一个具体的版本以保证构建的可复现性
FROM golang:1.23.3 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载所有依赖项。利用 Docker 的层缓存，只有在 go.mod/go.sum 变化时才会重新下载
RUN go mod download

# 复制项目的其余源代码
COPY . .

# 编译应用
# CGO_ENABLED=0: 禁用 CGO，以生成静态链接的二进制文件，这在 Alpine 这样的最小化镜像中很重要
# -o /app/xinde: 指定输出的二进制文件名为 xinde，并放在 /app 目录下
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/xinde ./cmd/app/main.go


# --- Stage 2: Production ---
# 使用一个极度轻量级的 Alpine 镜像作为最终的生产环境
FROM alpine:latest

# Alpine 默认没有 ca-certificates，对于发起 HTTPS 请求等操作是必需的
# tzdata 是为了让 Go 应用能正确处理时区
RUN apk add --no-cache ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件到当前阶段
COPY --from=builder /app/xinde .

# 复制配置文件目录
# 注意：这只是为了提供一个默认的配置文件，通常我们会通过挂载来覆盖它
COPY configs ./configs

# 暴露应用监听的端口 (请确保与你的 config.yaml 中配置的端口一致)
EXPOSE 8080

# 定义容器启动时执行的命令
# ["./xinde", "-f", "./configs/config.yaml"]
# ./xinde: 运行我们编译的程序
# -f ./configs/config.yaml: 使用我们改造后的 -f 参数，指定配置文件的路径
ENTRYPOINT ["./xinde"]
CMD ["-f", "./configs/config.yaml"]