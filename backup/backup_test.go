package backup

import "testing"

import "github.com/andreacioni/motionctrl/config"

func TestBackupCron(t *testing.T) {
	Init(config.Backup{When: "0 0 * * * *", Method: "gdrive"}, ".")
	Shutdown()
}

func TestBackupSize(t *testing.T) {
	Init(config.Backup{When: "10MB", Method: "gdrive"}, ".")
	Shutdown()
}
