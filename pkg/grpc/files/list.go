package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (f *FileServer) List(ctx context.Context, req *emptypb.Empty) (*ListResponse, error) {
	files, err := f.fm.List()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &ListResponse{Files: files}, nil
}
