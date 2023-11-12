package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	echo "github.com/EmptyLungs/golang_file_server/pkg/grpc/echo"
)

type Server struct {
	logger *zap.Logger
	config *Config
}

type Config struct {
	Port        int    `mapstructure:"grpc-port"`
	ServiceName string `mapstructure:"grpc-service-name"`
}

func NewServer(config *Config, logger *zap.Logger) (*Server, error) {
	childLogger := logger.With(zap.String("source", "grpc"))
	srv := &Server{
		logger: childLogger,
		config: config,
	}
	return srv, nil
}

func (s *Server) ListenAndServe() *grpc.Server {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.config.Port))
	if err != nil {
		s.logger.Fatal("Failed to start TCP listener", zap.Error(err))
	}
	s.logger.Info("Starting grpc server")
	srv := grpc.NewServer()
	server := health.NewServer()
	reflection.Register(srv)
	grpc_health_v1.RegisterHealthServer(srv, server)
	server.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	echoService := echo.NewEchoService(s.logger)
	echo.RegisterEchoServiceServer(srv, echoService)
	go func() {
		if err := srv.Serve(listener); err != nil {
			s.logger.Fatal("Failed to start grpc server", zap.Error(err))
		}
	}()
	return srv
}
