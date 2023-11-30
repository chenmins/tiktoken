# 使用一个基础镜像，例如Alpine Linux，因为它很小而且安全
FROM alpine:latest

# 将二进制文件复制到容器中
COPY tiktoken /usr/local/bin/

# 假设二进制文件名为ddns，设置二进制文件为容器的入口点
ENTRYPOINT ["/usr/local/bin/tiktoken"]
