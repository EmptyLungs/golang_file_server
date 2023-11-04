package files

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"

	"go.uber.org/zap"
)

type FileManager struct {
	workDir string
	logger  *zap.Logger
}

func NewFileManager(workDir string, logger *zap.Logger) (*FileManager, error) {
	err := checkDirExists(workDir)
	if err != nil {
		return nil, err
	}
	childLogger := logger.With(zap.String("source", "file_uploader"))
	fm := &FileManager{
		workDir: workDir,
		logger:  childLogger,
	}
	return fm, nil
}

func (fm *FileManager) CreateFile(file multipart.File, handler *multipart.FileHeader) error {
	if handler.Size == 0 {
		return errors.New("file is empty")
	}
	if handler.Size > 5*1024*1024 {
		return errors.New("file is too big")
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(fm.workDir, handler.Filename))
	if err != nil {
		return err
	}
	f.Write(fileBytes)
	fm.logger.Info("upload", zap.String("filename", handler.Filename), zap.Int64("size_bytes", handler.Size))
	return nil
}

func (fm *FileManager) DeleteFile() error {
	return nil
}

func (fm *FileManager) ListFiles() ([]string, error) {
	var fileNames []string
	files, err := os.ReadDir(fm.workDir)
	if err != nil {
		return nil, err
	}
	for _, e := range files {
		fileNames = append(fileNames, e.Name())
	}
	return fileNames, nil
}
