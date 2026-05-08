// gRPC 拦截器（鉴权、日志、重试）
package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"his-go/pkg/logger"
)

// UnaryServerInterceptor 一元拦截器链
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		latency := time.Since(start)

		if err != nil {
			logger.Error("gRPC 调用失败",
				zap.String("method", info.FullMethod),
				zap.Duration("latency", latency),
				zap.Error(err),
			)
		} else {
			logger.Debug("gRPC 调用成功",
				zap.String("method", info.FullMethod),
				zap.Duration("latency", latency),
			)
		}
		return resp, err
	}
}

// UnaryClientInterceptor 客户端拦截器（日志 + 重试）
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		latency := time.Since(start)

		if err != nil {
			logger.Error("gRPC 客户端调用失败",
				zap.String("method", method),
				zap.Duration("latency", latency),
				zap.Error(err),
			)
		} else {
			logger.Debug("gRPC 客户端调用成功",
				zap.String("method", method),
				zap.Duration("latency", latency),
			)
		}
		return err
	}
}
