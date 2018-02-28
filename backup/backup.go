package backup

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/andreacioni/aescrypt"
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

type State string

const (
	StateActiveIdle    State = "ACTIVE_IDLE"
	StateActiveRunning State = "ACTIVE_RUNNING"
	StateActiveErrored State = "ACTIVE_ERRORED"
	StateDeactivated   State = "NOT_ACTIVE"
)

const (
	GoogleDriveMethod = "google"
	TestMockMethod    = "mock"
)

var (
	backupConfig    config.Backup
	targetDirectory string

	dirWatcher *watcher.Watcher

	cronSheduler *cron.Cron

	workerMutex trylock.Mutex
	runFlag     = abool.New()

	uploadService UploadService

	backupStatus     State
	backupStateMutex sync.Mutex
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

	setStatus(StateActiveIdle)
	runFlag.Set()
}

func GetStatus() State {
	backupStateMutex.Lock()
	defer backupStateMutex.Unlock()
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

	setStatus(StateDeactivated)

}

func setStatus(status State) {
	backupStateMutex.Lock()
	defer backupStateMutex.Unlock()
	backupStatus = status
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
		glg.Infof("Setup directory watcher on: %s, max size: %d bytes", targetDirectory, maxFolderSize)

		dirWatcher = watcher.New()

		dirWatcher.SetMaxEvents(1)

		dirWatcher.FilterOps(watcher.Create)

		if err := dirWatcher.AddRecursive(targetDirectory); err != nil {
			return err
		}

		go func() {
			for {
				select {
				case <-dirWatcher.Event:
					folderSize := evaluateFolderSize()
					glg.Debugf("Output folder size: %d/%d", folderSize, maxFolderSize)
					if folderSize > maxFolderSize {
						go backupWorker()
					}
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
		return &GoogleDriveBackupService{}, nil
	default:
		err = fmt.Errorf("Backup method not recognized")
	}
	return nil, err
}

func listFile(targetDirectory string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(targetDirectory, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() {
				fullPath, err := filepath.Abs(path)
				if err != nil {
					glg.Error(err)
				} else {
					fileList = append(fileList, fullPath)
				}

			}
		}

		return err
	})

	return fileList, err

}

func archiveFiles(fileList []string) (string, error) {
	glg.Debugf("Now compressing: %+v", fileList)

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, f := range fileList {
		fileInfo, err := os.Stat(f)
		if err != nil {
			return "", err
		}
		hdr := &tar.Header{
			Name: filepath.Base(f),
			Mode: 0600,
			Size: fileInfo.Size(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return "", err
		}
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return "", err
		}
		if _, err := tw.Write(b); err != nil {
			return "", err
		}
	}
	if err := tw.Close(); err != nil {
		return "", err
	}

	archiveFileName := time.Now().Format("20060102_150405") + ".tar.gz"
	archiveFilePath := filepath.Join(os.TempDir(), archiveFileName)
	if err := ioutil.WriteFile(archiveFilePath, buf.Bytes(), 0600); err != nil {
		return "", err
	}

	if err := removeFiles(fileList); err != nil {
		return "", err
	}

	glg.Debugf("Archive file created:%s", archiveFileName)

	return archiveFilePath, nil
}

func encryptAndUpload(filePath string, key string) error {
	var err error

	if key != "" {
		glg.Debugf("Encryption enabled for file: %s", filePath)

		aesFilePath := filePath + ".aes"
		if err = aescrypt.New(key).Encrypt(filePath, aesFilePath); err != nil {
			return fmt.Errorf("Failed to encrypt file %s: %v", filePath, aesFilePath)
		}
		filePath = aesFilePath
	}

	err = uploadService.Upload(filePath)

	if err != nil {
		return err
	}

	glg.Debugf("Uploaded file: %s", filePath)

	return nil
}

func removeFiles(filePath []string) error {
	for _, f := range filePath {
		err := os.Remove(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func backupWorker() {
	if runFlag.IsSet() {
		if workerMutex.TryLock() {
			defer workerMutex.Unlock()
			glg.Debug("Backup service worker is running now")
			setStatus(StateActiveRunning)

			fileList, err := listFile(targetDirectory)

			glg.Debugf("Backup file list: %+v", fileList)

			if err != nil {
				glg.Error(err)
				setStatus(StateActiveErrored)
			} else {
				if backupConfig.Archive { //Group file into archive, encrypt (if needed) and upload
					utils.BlockSlideSlice(fileList, backupConfig.FilePerArchive, func(subList interface{}) bool {
						subFileList := subList.([]string)

						archive, err := archiveFiles(subFileList)

						if err != nil {
							return false
						}

						err = encryptAndUpload(archive, backupConfig.EncryptionKey)

						if err != nil {
							return false
						}

						if !backupConfig.KeepLocalCopy {
							subFileList = append(subFileList, archive) //Remove archive file and photo
						}
						err = removeFiles(subFileList)

						if err != nil {
							return false
						}

						if !runFlag.IsSet() {
							return false
						}

						return true
					})

					if err != nil {
						glg.Error(err)
						setStatus(StateActiveErrored)
					}

					if !runFlag.IsSet() {
						return
					}
				} else { //Encrypt file (if needed) and upload. No archive
					for _, f := range fileList {

						err = encryptAndUpload(f, backupConfig.EncryptionKey)

						if err != nil {
							glg.Error(err)
							setStatus(StateActiveErrored)
							return
						}

						if !backupConfig.KeepLocalCopy {
							err = os.Remove(f)

							if err != nil {
								glg.Error(err)
								setStatus(StateActiveErrored)
								return
							}
						}

						if !runFlag.IsSet() {
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
