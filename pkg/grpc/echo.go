package grpc

import (
	"context"

	"go.uber.org/zap"
)

type EchoService struct {
	logger *zap.Logger
	UnimplementedEchoServiceServer
}

func (s EchoService) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	message := req.GetMessage()
	s.logger.Info("handling echo")
	response := &EchoResponse{
		Message: message,
	}
	return response, nil
}
