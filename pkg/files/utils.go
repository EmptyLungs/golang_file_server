package files

import "os"

func checkDirExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
		return err
	}
	return nil
}
