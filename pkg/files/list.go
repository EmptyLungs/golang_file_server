package files

import "io/fs"

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
