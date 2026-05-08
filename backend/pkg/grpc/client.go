// Package grpc gRPC 客户端连接池封装
package grpc

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second,
	Timeout:             3 * time.Second,
	PermitWithoutStream: true,
}

// NewGrpcClient 创建 gRPC 客户端连接
func NewGrpcClient(host string, port int) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 gRPC 客户端连接失败: %w", err)
	}

	return conn, nil
}

// NewGrpcServer 创建 gRPC 服务端
func NewGrpcServer() *grpc.Server {
	return grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024),
		grpc.MaxSendMsgSize(10*1024*1024),
	)
}
