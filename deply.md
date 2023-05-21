第一行指定了基础镜像为 golang:1.16，这是一个预安装了 Go 语言环境的镜像。

WORKDIR /app 将容器的工作目录设置为 /app。

COPY main.go go.* /app/ 将当前目录下的 main.go 文件和 go.mod、go.sum 文件复制到容器的 /app 目录中。

RUN go mod download && go build -o main . 在容器中执行两个命令。首先执行 go mod download 命令下载项目的依赖包，然后执行 go build -o main . 命令构建可执行文件 main。

CMD ["./main"] 设置容器启动时的命令，即运行可执行文件 main。

保存并关闭 Dockerfile 文件。

使用以下命令在 Docker 中构建镜像：

docker build -t your-image-name .