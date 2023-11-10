package files

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockFSWrapper struct {
	fs  fstest.MapFS
	err error
}

func (m MockFSWrapper) ReadDir(name string) ([]fs.DirEntry, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.fs.ReadDir(name)
}

func (m MockFSWrapper) Open(name string) (fs.File, error) {
	return m.fs.Open(name)
}

func TestFileManagerList(t *testing.T) {
	assert := assert.New(t)
	mfs := MockFSWrapper{
		err: nil,
		fs: fstest.MapFS{
			"test.txt":          &fstest.MapFile{},
			"pleaseignore.json": &fstest.MapFile{},
		},
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fileManager := FileManager{fs: mfs, workDir: ".", logger: logger}
	files, err := fileManager.List()
	assert.Nil(err)
	assert.ElementsMatch([]string{"test.txt", "pleaseignore.json"}, files)
}
func TestFileManagerList_Fail(t *testing.T) {
	assert := assert.New(t)

	errMsg := "Failed to list working directory"
	mfs := MockFSWrapper{
		err: errors.New(errMsg),
		fs:  fstest.MapFS{},
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fileManager := FileManager{fs: mfs, workDir: ".", logger: logger}
	files, err := fileManager.List()
	t.Logf(err.Error())
	assert.Nil(files)
	assert.Equal(errMsg, err.Error())
}
