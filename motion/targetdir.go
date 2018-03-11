package motion

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/andreacioni/motionctrl/utils"
)

type TargetDirFile struct {
	Name         string    `json:"filename"`
	CreationTime time.Time `json:"creationDate"`
}

func TargetDirListFiles() ([]TargetDirFile, error) {
	fileInfo, _, _, err := utils.ListFiles(readOnlyConfig[ConfigTargetDir], nil)

	var files []TargetDirFile

	for _, f := range fileInfo {
		files = append(files, TargetDirFile{Name: f.Name(), CreationTime: f.ModTime()})
	}

	return files, err
}

func TargetDirSize() (int64, error) {
	_, _, size, err := utils.ListFiles(readOnlyConfig[ConfigTargetDir], nil)

	if err != nil {
		return -1, fmt.Errorf("Unable to evaluate size of target directory: %v", err)
	}

	return size, err
}

func TargetDirGetFile(filename string) (string, error) {
	return filepath.Join(readOnlyConfig[ConfigTargetDir], filename), nil
}

func TargetDirRemoveFile(filename string) error {
	if err := os.Remove(filepath.Join(readOnlyConfig[ConfigTargetDir], filename)); err != nil {
		return fmt.Errorf("Unable to remove %s: %v", filename, err)
	}

	return nil
}
