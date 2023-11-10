package files

import (
	"io"
	"os"
	"path"

	"go.uber.org/zap"
)

func (fm FileManager) Create(file io.Reader, filename string) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(fm.workDir, filename))
	if err != nil {
		return err
	}
	if _, err = f.Write(fileBytes); err != nil {
		return err
	}
	fm.logger.Info("upload", zap.String("filename", filename), zap.Int("size_bytes", len(fileBytes)))
	return nil
}
