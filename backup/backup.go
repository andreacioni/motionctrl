package backup

import (
	"log"
	"time"

	"github.com/kpango/glg"
	"github.com/radovskyb/watcher"

	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/motion"
)

var (
	dirWatcher   *watcher.Watcher
	backupConfig []config.Backup
)

func Init() {
	backupConfig = config.Get().Backup

	if backupConfig != nil && len(backupConfig) > 0 {
		setupWatcher()

		/*for _, _ := range backupConfig {
			val.
		}*/
	}

}

func setupWatcher() {
	dirWatcher = watcher.New()

	dirWatcher.SetMaxEvents(1)

	dirWatcher.FilterOps(watcher.Create)

	if err := dirWatcher.AddRecursive(motion.ConfigGetRO(motion.ConfigTargetDir)); err != nil {
		glg.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-dirWatcher.Event:
				glg.Debug(event)
			case err := <-dirWatcher.Error:
				glg.Fatal(err)
			case <-dirWatcher.Closed:
				return
			}
		}
	}()

	if err := dirWatcher.Start(time.Minute); err != nil {
		log.Fatalln(err)
	}
}

func evaluateFolderSize() int64 {
	var totSize int64
	for _, fInfo := range dirWatcher.WatchedFiles() {
		if !fInfo.IsDir() {
			totSize += fInfo.Size()
		}
	}

	return totSize
}

func buildBackupWorker() {

}
