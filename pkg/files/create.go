package files

import (
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func (fm FileManager) Create(file io.Reader, filename string) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(fm.workDir, filename))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(fileBytes); err != nil {
		return err
	}
	fm.logger.Info("upload", zap.String("filename", filename), zap.Int("size_bytes", len(fileBytes)))
	return nil
}
