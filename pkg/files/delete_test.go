package files

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestFileManagerDeleteFileNotExist(t *testing.T) {
	assert := assert.New(t)
	tempDir := t.TempDir()

	fs := os.DirFS(tempDir)
	logger := zap.NewNop()
	fileManager, err := NewFileManager(tempDir, fs, logger)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = fileManager.Delete("test.txt")
	assert.Error(err)
	assert.True(errors.Is(err, os.ErrNotExist))
}
