package grpc

import (
	"bytes"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileServer) Upload(stream FileService_UploadServer) error {
	// todo implement chunk-by-chunk write
	// to prevent accumulating all file in memory
	var fileSize uint32
	fileSize = 0
	var fileData []byte
	var fileName string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, "faield to read stream")
		}
		if req.GetFilename() != "" {
			fileName = req.GetFilename()
		}
		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		fileData = append(fileData, chunk...)
	}
	buffer := bytes.NewBuffer(fileData)
	if fileName == "" {
		return status.Error(codes.InvalidArgument, "missing file name")
	}
	if err := f.fm.Create(buffer, fileName); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return stream.SendAndClose(&UploadResponse{Message: "ok", Filename: fileName, Filesize: fileSize})
}
