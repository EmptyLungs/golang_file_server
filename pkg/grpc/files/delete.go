package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	filename := req.GetFilename()
	if err := f.fm.Delete(filename); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &DeleteResponse{Message: "ok"}, nil
}
