package files

import (
	"io"
	"io/fs"

	"go.uber.org/zap"
)

type IFileManager interface {
	Create(file io.Reader, filename string) error
	Delete(filename string) error
	List() ([]string, error)
}

type FileManager struct {
	workDir string
	fs      fs.FS
	logger  *zap.Logger
}

func NewFileManager(workDir string, fs fs.FS, logger *zap.Logger) (*FileManager, error) {
	err := checkDirExists(workDir)
	if err != nil {
		return nil, err
	}
	childLogger := logger.With(zap.String("source", "file_uploader"))
	fm := &FileManager{
		workDir: workDir,
		fs:      fs,
		logger:  childLogger,
	}
	return fm, nil
}
