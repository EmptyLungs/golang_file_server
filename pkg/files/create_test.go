package files

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFileManagerCreate(t *testing.T) {
	assert := assert.New(t)
	tempDir := t.TempDir()
	fs := os.DirFS(tempDir)

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fileManager, err := NewFileManager(tempDir, fs, logger)
	if err != nil {
		t.Fatalf(err.Error())
	}
	data := []byte("Hello world!")
	reader := bytes.NewReader(data)
	expectedFileName := "test.txt"
	expectedFilePath := filepath.Join(tempDir, expectedFileName)
	err = fileManager.Create(reader, expectedFileName)
	assert.Nil(err)

	content, err := os.ReadFile(expectedFilePath)
	assert.Nil(err)
	assert.Equal("Hello world!", string(content))
}
