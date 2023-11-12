package grpc

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	files "github.com/EmptyLungs/golang_file_server/pkg/files"
	grpc_echo "github.com/EmptyLungs/golang_file_server/pkg/grpc/echo"
	grpc_files "github.com/EmptyLungs/golang_file_server/pkg/grpc/files"
)

type Server struct {
	logger *zap.Logger
	config *Config
	fm     files.IFileManager
}

type Config struct {
	Port        int    `mapstructure:"grpc-port"`
	ServiceName string `mapstructure:"grpc-service-name"`
}

func NewServer(config *Config, logger *zap.Logger, fileManager files.IFileManager) (*Server, error) {
	childLogger := logger.With(zap.String("source", "grpc"))
	srv := &Server{
		logger: childLogger,
		config: config,
		fm:     fileManager,
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

	echoService := grpc_echo.NewEchoService(s.logger)
	grpc_echo.RegisterEchoServiceServer(srv, echoService)

	fileService := grpc_files.NewService(s.logger, s.fm)
	grpc_files.RegisterFileServiceServer(srv, fileService)
	go func() {
		if err := srv.Serve(listener); err != nil {
			s.logger.Fatal("Failed to start grpc server", zap.Error(err))
		}
	}()
	return srv
}
