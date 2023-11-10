package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDirExists(t *testing.T) {
	assert := assert.New(t)
	tmp := t.TempDir()
	err := checkDirExists(tmp)
	assert.Nil(err)
}

func TestCheckDirNotExists(t *testing.T) {
	assert := assert.New(t)
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test")
	err := checkDirExists(path)
	assert.Nil(err)
	info, err := os.Stat(path)
	assert.Nil(err)
	assert.Equal("test", info.Name())
}

// todo: refactor checkDirExists
// redo os.Stat with fs.FS
// rename func
