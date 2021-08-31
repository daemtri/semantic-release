# 编译后段代码，前端代码可以使用 go embed功能打包进go二进制，也可以在下一步copy进最终镜像
FROM harbor.bianfeng.com/library/golang:1.17-alpine AS Builder
ARG TAGS="timetzdata"

WORKDIR /src

ENV CGO_ENABLED=0
ENV GOPRIVATE="hub.imeete.com,git.imeete.com,git.bianfeng.com"
ENV GOPROXY="https://goproxy.io"

# 缓存
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN cd cmd/semantic-release && go build -tags "${TAGS}" -o /src/dist/semantic-release

FROM alpine

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
COPY --from=Builder --chown=0:0 /src/dist/semantic-release /usr/local/bin/semantic-release

ENTRYPOINT ["semantic-release"]
