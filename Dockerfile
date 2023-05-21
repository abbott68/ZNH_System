# 使用 golang 镜像作为基础镜像
FROM golang:1.16

# 设置工作目录
WORKDIR /app

# 将 Go 代码复制到容器中
COPY main.go go.* /app/

# 安装依赖并构建可执行文件
RUN go mod download && go build -o main .

# 设置容器启动命令
CMD ["./main"]
