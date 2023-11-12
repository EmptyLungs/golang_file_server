package grpc

import "context"

func (f *FileServer) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	filename := req.GetFilename()
	if err := f.fm.Delete(filename); err != nil {
		return nil, err
	}
	return &DeleteResponse{Message: "ok"}, nil
}
