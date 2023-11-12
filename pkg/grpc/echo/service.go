package grpc

import (
	"context"

	"go.uber.org/zap"
)

type EchoService struct {
	logger *zap.Logger
	UnimplementedEchoServiceServer
}

func NewEchoService(logger *zap.Logger) *EchoService {
	return &EchoService{
		logger: logger,
	}
}

func (s EchoService) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	message := req.GetMessage()
	s.logger.Info("handling echo")
	response := &EchoResponse{
		Message: message,
	}
	return response, nil
}
