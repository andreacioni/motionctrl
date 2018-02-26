package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/utils"

	"github.com/LK4D4/trylock"
	"github.com/kpango/glg"
	"github.com/radovskyb/watcher"
	"github.com/robfig/cron"
	"github.com/tevino/abool"
)

type UploadService interface {
	Authenticate() error
	Upload(string) error
}

type State int

const (
	StateActiveIdle State = iota
	StateActiveRunning
	StateActiveErrored
	StateDeactivated
)

const (
	GoogleDriveMethod = "google"
	TestMockMethod    = "mock"
)

var (
	backupConfig    config.Backup
	targetDirectory string

	dirWatcher    *watcher.Watcher
	maxFolderSize uint64

	cronSheduler *cron.Cron

	mu      trylock.Mutex
	runFlag = abool.New()

	uploadService UploadService

	backupStatus State
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

	uploadService, err = buildUploadService(backConf.Method)

	if err != nil {
		glg.Fatal(err)
	}

	err = uploadService.Authenticate()

	if err != nil {
		glg.Fatalf("Failed to authenticate on upload service: %s", err.Error())
	}

	backupStatus = StateActiveIdle
	runFlag.Set()
}

func Status() State {
	return backupStatus
}

func Shutdown() {

	glg.Info("Shuting down backup service")

	runFlag.UnSet()

	if dirWatcher != nil {
		dirWatcher.Close()
		dirWatcher = nil
	}

	if cronSheduler != nil {
		cronSheduler.Stop()
		cronSheduler = nil
	}

	backupStatus = StateDeactivated

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
			return err
		}

		go func() {
			for {
				select {
				case event := <-dirWatcher.Event:
					glg.Debug(event)
				case err := <-dirWatcher.Error:
					glg.Error(err)
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

func buildUploadService(uploadMethod string) (UploadService, error) {
	var err error
	switch uploadMethod {
	case GoogleDriveMethod:
		return &GoogleDriveBackupService{}, nil
	case TestMockMethod:
		return &MockBackupService{}, nil
	default:
		err = fmt.Errorf("Backup method not recognized")
	}
	return nil, err
}

func listFile(targetDirectory string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(targetDirectory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fullPath, err := filepath.Abs(path)
			if err != nil {
				glg.Error(err)
			} else {
				fileList = append(fileList, fullPath)
			}

		}
		return nil
	})

	return fileList, err

}

func archiveFiles(fileList []string) (string, error) { //TODO
	glg.Debugf("Now compressing: %+v", fileList)

	return "", nil
}

func encryptAndUpload(filePath string, key string) error {
	var err error

	if key != "" { //TODO

	}

	err = uploadService.Upload(filePath)

	return err
}

func backupWorker() {
	if runFlag.IsSet() {
		if mu.TryLock() {
			defer mu.Unlock()
			glg.Debug("Backup service worker is running now")
			backupStatus = StateActiveRunning

			fileList, err := listFile(targetDirectory)

			glg.Debugf("Backup file list: %+v", fileList)

			if err != nil {
				glg.Error(err)
				backupStatus = StateActiveErrored
			} else {
				if backupConfig.Archive {
					utils.BlockSlideSlice(fileList, backupConfig.FilePerArchive, func(subList interface{}) bool {
						subFileList := subList.([]string)

						zipArchive, err := archiveFiles(subFileList)

						if err != nil {
							return false
						}

						err = encryptAndUpload(zipArchive, backupConfig.Key)

						if err != nil || !runFlag.IsSet() {
							return false
						}

						return true
					})

					if err != nil {
						glg.Error(err)
						backupStatus = StateActiveErrored
					}

					if runFlag.IsSet() {
						return
					}
				} else {
					for _, f := range fileList {

						err = encryptAndUpload(f, backupConfig.Key)

						if err != nil {
							glg.Error(err)
							backupStatus = StateActiveErrored
							return
						}

						if runFlag.IsSet() {
							return
						}
					}
				}

			}

		} else {
			glg.Debug("Backup worker is already running")
		}
	} else {
		glg.Warn("Worker disabled")
	}

}
