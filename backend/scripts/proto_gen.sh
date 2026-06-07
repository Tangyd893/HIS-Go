#!/bin/bash
# Proto 代码生成脚本
# 用法: bash scripts/proto_gen.sh
# 依赖: protoc, protoc-gen-go, protoc-gen-go-grpc
#
# 安装依赖:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#   确保 $GOPATH/bin 在 PATH 中

set -e

PROTO_DIR="api/proto"
OUT_DIR="api/proto"

echo "=== 开始生成 Proto Go 代码 ==="

command -v protoc >/dev/null 2>&1 || { echo "错误: protoc 未安装，请执行: apt install protobuf-compiler 或从 https://github.com/protocolbuffers/protobuf/releases 下载"; exit 1; }
command -v protoc-gen-go >/dev/null 2>&1 || { echo "错误: protoc-gen-go 未安装，请执行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1; }
command -v protoc-gen-go-grpc >/dev/null 2>&1 || { echo "错误: protoc-gen-go-grpc 未安装，请执行: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; exit 1; }

# 先生成公共类型
echo "生成公共类型..."
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${OUT_DIR} \
    --go_opt=paths=source_relative \
    ${PROTO_DIR}/common/common.proto

# 生成各服务 proto
for dir in $(ls -d ${PROTO_DIR}/*/); do
    service=$(basename "$dir")
    if [[[ "$service" = "common" ]]; then
        continue
    fi
    echo "生成服务: ${service}"

    protoc \
        --proto_path=${PROTO_DIR} \
        --go_out=${OUT_DIR} \
        --go_opt=paths=source_relative \
        --go-grpc_out=${OUT_DIR} \
        --go-grpc_opt=paths=source_relative \
        ${PROTO_DIR}/${service}/*.proto
done

echo "=== Proto 代码生成完成 ==="
echo ""
echo "生成的文件位于 api/proto/ 各子目录下："
echo "  - *.pb.go       (消息定义)"
echo "  - *_grpc.pb.go   (gRPC 服务接口与注册函数)"
