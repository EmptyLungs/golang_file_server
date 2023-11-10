package files

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFileManagerCreate(t *testing.T) {
	assert := assert.New(t)
	tempDir := t.TempDir()
	fs := os.DirFS(tempDir)

	logger := zap.NewNop()
	fileManager, err := NewFileManager(tempDir, fs, logger)
	if err != nil {
		t.Fatalf(err.Error())
	}
	reader := strings.NewReader("Hello world!")
	expectedFileName := "test.txt"
	expectedFilePath := filepath.Join(tempDir, expectedFileName)
	err = fileManager.Create(reader, expectedFileName)
	assert.Nil(err)

	content, err := os.ReadFile(expectedFilePath)
	assert.Nil(err)
	assert.Equal("Hello world!", string(content))
}

func TestFileManagerCreate_BadReaderFail(t *testing.T) {
	assert := assert.New(t)
	tempDir := t.TempDir()
	fs := os.DirFS(tempDir)

	logger := zap.NewNop()
	fileManager, err := NewFileManager(tempDir, fs, logger)
	if err != nil {
		t.Fatalf(err.Error())
	}
	excpectedErr := errors.New("Read file error")
	badFile := iotest.ErrReader(excpectedErr)

	err = fileManager.Create(badFile, "test.json")
	assert.NotNil(err)
	assert.Equal(excpectedErr, err)
}
