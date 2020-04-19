FROM golang:1.14-alpine as builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /working
COPY .. /working
EXPOSE 8080
# 更新安装源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add --no-cache bash supervisor curl strace gdb

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin-blog-api main.go


Copy ./supervisord.conf /etc/supervisord.conf


CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]