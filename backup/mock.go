package backup

import (
	"github.com/andreacioni/motionctrl/config"
)

type MockBackupService struct {
}

func (b *MockBackupService) Upload(file string) error {
	return nil
}

func (b *MockBackupService) Setup(backConf config.Backup) error {
	return nil
}
