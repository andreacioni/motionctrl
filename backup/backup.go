package backup

import (
	"os"
	"time"

	"github.com/andreacioni/motionctrl/utils"

	"github.com/kpango/glg"
	"github.com/radovskyb/watcher"
	"github.com/robfig/cron"

	"github.com/andreacioni/motionctrl/config"

	"github.com/LK4D4/trylock"
)

type UploadService interface {
	Upload(os.File) error
}

var (
	backupConfig    config.Backup
	targetDirectory string

	dirWatcher    *watcher.Watcher
	maxFolderSize uint64

	cronSheduler *cron.Cron

	mu trylock.Mutex

	uploadService UploadService
)

func Init(backConf config.Backup, targetDir string) {
	backupConfig = backConf
	targetDirectory = targetDir

	glg.Debugf("Initializing backup service: %+v, target directory: %s", backConf, targetDir)

	err := setupCronSheduler()

	if err != nil {
		err := setupDirectoryWatcher()

		if err != nil {
			glg.Fatalf("Not a valid size/cron expression in' backup.when'=%s", backupConfig.When)
		}
	}

	uploadService = buildUploadService(backConf.Method)

}

func Shutdown() {

	glg.Info("Shuting down backup service")

	if dirWatcher != nil {
		dirWatcher.Close()
		dirWatcher = nil
	}

	if cronSheduler != nil {
		cronSheduler.Stop()
		cronSheduler = nil
	}
}

func setupCronSheduler() error {
	cronSheduler = cron.New()

	err := cronSheduler.AddFunc(backupConfig.When, backupWorker)

	if err == nil {
		glg.Infof("Cron job is started, running on: %s", backupConfig.When)
		cronSheduler.Start()
	} else {
		cronSheduler = nil
	}

	return err
}

func setupDirectoryWatcher() error {

	maxFolderSize, err := utils.FromTextSize(backupConfig.When)

	if err == nil {
		glg.Infof("Setup directory watcher, max size: %d bytes", maxFolderSize)

		dirWatcher = watcher.New()

		dirWatcher.SetMaxEvents(1)

		dirWatcher.FilterOps(watcher.Create)

		if err := dirWatcher.AddRecursive(targetDirectory); err != nil {
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

		go func() {
			//Running directory watcher every 1 minute
			if err := dirWatcher.Start(time.Minute); err != nil {
				glg.Error(err)
			}
		}()

		dirWatcher.Wait()
	}

	return err

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

func buildUploadService(uploadMethod string) UploadService {

}

func backupWorker() {
	if mu.TryLock() {
		defer mu.Unlock()
		glg.Debug("Backup service worker is running now")

	} else {
		glg.Debug("Backup worker is already running")
	}
}
