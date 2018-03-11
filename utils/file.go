package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func ListFiles(dir string, filterFunc func(os.FileInfo) bool) ([]os.FileInfo, []string, int64, error) {
	fileList := []string{}
	fileInfo := []os.FileInfo{}
	var folderSize int64

	fInfo, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, nil, 0, err
	}

	//We get only regular files (no directory, symbolic files, hidden files) and older than 30 sec
	for _, f := range fInfo {
		if filterFunc == nil || filterFunc(f) {
			filePath, err := filepath.Abs(filepath.Join(dir, f.Name()))

			if err != nil {
				return nil, nil, 0, err
			}
			fileList = append(fileList, filePath)
			fileInfo = append(fileInfo, f)
			folderSize += f.Size()
		}
	}

	return fileInfo, fileList, folderSize, err

}
