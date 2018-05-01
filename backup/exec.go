package backup

import (
	"fmt"
	"github.com/andreacioni/motionctrl/config"
	"os/exec"
)

type ExecBackupService struct {
	execCmd string
}

func (b *ExecBackupService) Upload(file string) error {
	cmd := exec.Command(b.execCmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to execute backup command %s, %v", b.execCmd, err)
	}

	return nil
}

func (b *ExecBackupService) Setup(backConf config.Backup) error {
	b.execCmd = backConf.Command
	return nil
}
