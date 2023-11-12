package grpc

import (
	"github.com/EmptyLungs/golang_file_server/pkg/files"
	"go.uber.org/zap"
)

type FileServer struct {
	UnimplementedFileServiceServer
	logger *zap.Logger
	fm     files.IFileManager
}

func NewService(logger *zap.Logger, filemanager files.IFileManager) *FileServer {
	return &FileServer{
		logger: logger,
		fm:     filemanager,
	}
}
