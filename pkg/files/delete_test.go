package files

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func TestFileManagerDelete(t *testing.T) {
	tempDir := t.TempDir()

	fs := os.DirFS(tempDir)
	logger := zap.NewNop()
	fileManager, err := NewFileManager(tempDir, fs, logger)
	filename := "test.txt"
	if err != nil {
		t.Fatalf(err.Error())
	}
	f, err := os.Create(filepath.Join(tempDir, filename))
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer f.Close()
	if _, err := f.Write([]byte("Hello world!")); err != nil {
		t.Fatalf(err.Error())
	}
	if err = fileManager.Delete(filename); err != nil {
		t.Fatalf(err.Error())
	}
}
