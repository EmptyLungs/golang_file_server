package grpc

import (
	"bytes"
	"errors"
	"io"
)

func (f *FileServer) Upload(stream FileService_UploadServer) error {
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
			return errors.New("error")
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
		return errors.New("missing filename")
	}
	if err := f.fm.Create(buffer, fileName); err != nil {
		return err
	}
	return stream.SendAndClose(&UploadResponse{Message: "ok"})
}
