package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"app/pkg/log"
)

type Server struct {
	*grpc.Server
	host   string
	port   int
	logger *log.Logger
}

type Option func(s *Server)

func NewServer(logger *log.Logger, opts ...Option) *Server {
	kaParams := keepalive.ServerParameters{
		Time:    10 * time.Second, // 10秒内没有请求则发送ping包
		Timeout: 2 * time.Second,  // 等待2秒超时
	}
	s := &Server{
		Server: grpc.NewServer(grpc.KeepaliveParams(kaParams)),
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		s.logger.Sugar().Fatalf("Failed to listen: %v", err)
	}
	if err = s.Server.Serve(lis); err != nil {
		s.logger.Sugar().Fatalf("Failed to serve: %v", err)
	}
	return nil

}
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	s.Server.GracefulStop()

	s.logger.Info("Server exiting")

	return nil
}
