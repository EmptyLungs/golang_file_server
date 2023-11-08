package files

import (
	"io"
	"io/fs"
	"os"
	"path"

	"go.uber.org/zap"
)

type IFileManager interface {
	Create(file io.Reader, filename string) error
	Delete() error
	List() ([]string, error)
}

type FileManager struct {
	workDir string
	fs      fs.FS
	logger  *zap.Logger
}

func NewFileManager(workDir string, logger *zap.Logger) (*FileManager, error) {
	err := checkDirExists(workDir)
	if err != nil {
		return nil, err
	}
	childLogger := logger.With(zap.String("source", "file_uploader"))
	fs := os.DirFS(workDir)
	fm := &FileManager{
		workDir: workDir,
		fs:      fs,
		logger:  childLogger,
	}
	return fm, nil
}

func (fm FileManager) Create(file io.Reader, filename string) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(fm.workDir, filename))
	if err != nil {
		return err
	}
	f.Write(fileBytes)
	fm.logger.Info("upload", zap.String("filename", filename), zap.Int("size_bytes", len(fileBytes)))
	return nil
}

func (fm FileManager) Delete() error {
	return nil
}

func (fm FileManager) List() ([]string, error) {
	var fileNames []string
	files, err := fs.ReadDir(fm.fs, ".")
	if err != nil {
		return nil, err
	}
	for _, e := range files {
		fileNames = append(fileNames, e.Name())
	}
	return fileNames, nil
}
