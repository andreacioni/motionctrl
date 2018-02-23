package backup

import (
	"os"
)

type GoogleDriveBackupService struct {
}

func (b GoogleDriveBackupService) Upload(f os.FileInfo) error {
	return nil
}
