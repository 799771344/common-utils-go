package main

import (
	"context"

	"google.golang.org/grpc"
)

// MiddlewareFunc 是一个 gRPC 中间件函数类型
type MiddlewareFunc func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)

// UnaryServerInterceptor 是一个 gRPC 服务器拦截器类型
type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)

// RegisterUnaryServerMiddleware 注册一个 gRPC 服务器的中间件
func RegisterUnaryServerMiddleware(srv *grpc.Server, middleware MiddlewareFunc) {
	// 保存原有的拦截器函数
	oldInterceptor := grpc.UnaryInterceptor

	// 定义新的拦截器函数
	newInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 调用中间件函数
		return middleware(ctx, req, info, handler)
	}

	// 将新的拦截器函数设置为 gRPC 服务器的拦截器
	grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 如果原有的拦截器函数不为空，则先调用原有的拦截器函数
		if oldInterceptor != nil {
			return oldInterceptor(ctx, req, info, newInterceptor)
		}

		// 否则直接调用新的拦截器函数
		return newInterceptor(ctx, req, info, handler)
	})
}
