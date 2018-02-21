package backup

import (
	"github.com/kpango/glg"
	"github.com/radovskyb/watcher"

	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/motion"
)

var (
	dirWatcher   watcher.Watcher
	backupConfig []motion.Backup
)

func Init() {
	backupConfig = config.Get().Backup

	if backupConfig && len(backupConfig) > 0 {
		setupWatcher()

		for _, val := range backupConfig {

		}
	}

}

func setupWatcher() {
	dirWatcher := watcher.New()

	if err := dirWatcher.AddRecursive(motion.ConfigGet(motion.TargetDir)); err != nil {
		glg.Fatal(err)
	}
}
