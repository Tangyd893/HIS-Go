#!/bin/bash
# Proto 代码生成脚本
# 用法: bash scripts/proto_gen.sh
# 依赖: protoc, protoc-gen-go, protoc-gen-go-grpc

set -e

PROTO_DIR="api/proto"
OUT_DIR="api/proto"

echo "=== 开始生成 Proto Go 代码 ==="

# 检查依赖
command -v protoc >/dev/null 2>&1 || { echo "错误: protoc 未安装"; exit 1; }
command -v protoc-gen-go >/dev/null 2>&1 || { echo "错误: protoc-gen-go 未安装，请执行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1; }
command -v protoc-gen-go-grpc >/dev/null 2>&1 || { echo "错误: protoc-gen-go-grpc 未安装，请执行: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; exit 1; }

for dir in $(ls -d ${PROTO_DIR}/*/); do
    service=$(basename "$dir")
    echo "处理服务: ${service}"

    protoc \
        --proto_path=${PROTO_DIR} \
        --go_out=${OUT_DIR} \
        --go_opt=paths=source_relative \
        --go-grpc_out=${OUT_DIR} \
        --go-grpc_opt=paths=source_relative \
        ${PROTO_DIR}/${service}/*.proto
done

echo "=== Proto 代码生成完成 ==="
