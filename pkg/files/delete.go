package files

import (
	"os"
	"path/filepath"
)

func (fm FileManager) Delete(path string) error {
	if err := os.Remove(filepath.Join(fm.workDir, path)); err != nil {
		return err
	}
	return nil
}
