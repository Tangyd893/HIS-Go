# 阶段一：编译
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git protobuf-dev

ARG SERVICE_NAME

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/${SERVICE_NAME}

# 阶段二：运行
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

EXPOSE 8080 8081 8082 8083 8084 8085 8086 8087 8088 8089 8090 8091 8092 8093 8094 8095 8096 8097
EXPOSE 9081 9082 9083 9084 9085 9086 9087 9088 9089 9090 9091 9092 9093 9094 9095 9096 9097

ENTRYPOINT ["./server"]
