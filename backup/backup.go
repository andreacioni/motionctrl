package backup

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/andreacioni/aescrypt"
	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/utils"
	"github.com/andreacioni/motionctrl/version"

	"github.com/LK4D4/trylock"
	"github.com/kpango/glg"
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
	StateDeactivated   State = "NOT_ACTIVE"
)

const (
	GoogleDriveMethod = "google"
	TestMockMethod    = "mock"
)

var (
	backupConfig    config.Backup
	targetDirectory string

	cronSheduler *cron.Cron

	maxFolderSize int64

	workerMutex trylock.Mutex
	runFlag     = abool.New()

	uploadService UploadService

	backupStatus     State
	backupStateMutex sync.Mutex
)

func Init(conf config.Backup, targetDir string) error {
	backupConfig = conf
	targetDirectory = targetDir

	glg.Debugf("Initializing backup service: %+v, target directory: %s", backupConfig, targetDir)

	uploadService, err := buildUploadService(backupConfig)

	if err != nil {
		glg.Warn("Backup service won't be active")
	} else {

		if err = uploadService.Authenticate(); err != nil {
			return fmt.Errorf("Failed to authenticate on upload service: %s", err.Error())
		}

		if err = setupCronSheduler(backupConfig); err != nil {

			if err = setupDirectoryWatcher(backupConfig, targetDir); err != nil {
				return fmt.Errorf("Not a valid size/cron expression in' backup.when'=%s", backupConfig.When)
			}
		}

		setStatus(StateActiveIdle)
		runFlag.Set()
	}

	return nil
}

func GetStatus() State {
	backupStateMutex.Lock()
	defer backupStateMutex.Unlock()
	return backupStatus
}

func Shutdown() {

	glg.Info("Shuting down backup service")

	runFlag.UnSet()

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

func setupCronSheduler(conf config.Backup) error {
	cronSheduler = cron.New()

	err := cronSheduler.AddFunc(conf.When, backupWorker)

	if err == nil {
		glg.Infof("Cron job is started, running on: %s", conf.When)
		cronSheduler.Start()
	} else {
		cronSheduler = nil
	}

	return err
}

func setupDirectoryWatcher(conf config.Backup, outDirectory string) error {

	maxFolderSize, err := utils.FromTextSize(conf.When)

	if err == nil {
		glg.Infof("Setup directory watcher on: %s, max size: %d bytes", outDirectory, maxFolderSize)

		cronSheduler = cron.New()

		if err = cronSheduler.AddFunc("@every 1m", checkSize); err == nil {
			glg.Info("Running directory size watcher every 1 minute")
			cronSheduler.Start()
		} else {
			cronSheduler = nil
		}
	}

	return err

}

func backuppableFile(f os.FileInfo) bool {
	return f != nil && f.Mode().IsRegular() && !strings.HasPrefix(f.Name(), ".") && f.ModTime().Add(time.Second*30).Before(time.Now())
}

func buildUploadService(conf config.Backup) (UploadService, error) {
	var err error
	switch conf.Method {
	case GoogleDriveMethod:
		return &GoogleDriveBackupService{}, nil
	case TestMockMethod:
		return &MockBackupService{}, nil
	default:
		err = fmt.Errorf("Backup method not found or invalid")
	}
	return nil, err
}

func listFile(targetDirectory string) ([]os.FileInfo, []string, int64, error) {
	fileList := []string{}
	fileInfo := []os.FileInfo{}
	var folderSize int64

	fInfo, err := ioutil.ReadDir(targetDirectory)

	if err != nil {
		return nil, nil, 0, err
	}

	//We get only regular files (no directory, symbolic files, hidden files) and older than 30 sec
	for _, f := range fInfo {
		if backuppableFile(f) {
			absPath, err := filepath.Abs(filepath.Join(targetDirectory, f.Name()))
			if err != nil {
				return nil, nil, 0, err
			}
			fileList = append(fileList, absPath)
			fileInfo = append(fileInfo, f)
			folderSize += f.Size()
		}
	}

	return fileInfo, fileList, folderSize, err

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

	archiveFileName := time.Now().Format(fmt.Sprintf("%s_20060102_150405", version.Name)) + ".tar.gz"
	archiveFilePath := filepath.Join(os.TempDir(), archiveFileName)
	if err := ioutil.WriteFile(archiveFilePath, buf.Bytes(), 0600); err != nil {
		return "", err
	}

	glg.Debugf("Archive file created:%s", archiveFilePath)

	return archiveFilePath, nil
}

func encryptAndUpload(filePath string, key string) (string, error) {
	var err error

	if key != "" {
		glg.Debugf("Encryption enabled for file: %s", filePath)

		aesFilePath := filePath + ".aes"
		if err = aescrypt.New(key).Encrypt(filePath, aesFilePath); err != nil {
			return "", fmt.Errorf("Failed to encrypt file %s: %v", filePath, aesFilePath)
		}

		glg.Debugf("Removing unencrypted file: %s", filePath)

		if err = os.Remove(filePath); err != nil { //Remove original unencrypted file
			return "", fmt.Errorf("Failed to remove original file %s: %v", filePath, aesFilePath)
		}

		filePath = aesFilePath
	}

	glg.Debugf("Uploading file: %s", filePath)

	err = uploadService.Upload(filePath)

	if err != nil {
		return "", err
	}

	glg.Debugf("Uploaded")

	return filePath, nil
}

func removeFiles(filePath []string) error {
	for _, f := range filePath {
		glg.Debugf("Removing file: %s", f)
		err := os.Remove(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkSize() {
	_, _, folderSize, err := listFile(targetDirectory)
	if err != nil {
		glg.Errorf("Failed to evaluate folder size: %v", err)
	}

	glg.Debugf("Output folder size: %d/%d", folderSize, maxFolderSize)
	if folderSize > maxFolderSize {
		go backupWorker()
	}
}

func backupWorker() {
	if runFlag.IsSet() {
		if workerMutex.TryLock() {
			defer workerMutex.Unlock()
			glg.Debug("Backup service worker is running now")
			setStatus(StateActiveRunning)

			_, fileList, _, err := listFile(targetDirectory)

			glg.Debugf("Backup file list: %+v", fileList)

			if err != nil {
				glg.Error(err)
			} else {
				if backupConfig.Archive { //Group file into archive, encrypt (if needed) and upload
					utils.BlockSlideSlice(fileList, backupConfig.FilePerArchive, func(subList interface{}) bool {
						subFileList := subList.([]string)

						archive, err := archiveFiles(subFileList)

						if err != nil {
							return false
						}

						if archive, err = encryptAndUpload(archive, backupConfig.EncryptionKey); err != nil {
							return false
						}

						subFileList = append(subFileList, archive) //Remove archive file and photo

						if err = removeFiles(subFileList); err != nil {
							return false
						}

						if !runFlag.IsSet() {
							glg.Warn("Run flag disabled")
							return false
						}

						return true
					})

					if err != nil {
						glg.Error(err)
						return
					}

					if !runFlag.IsSet() {
						glg.Warn("Run flag disabled")
						return
					}
				} else { //Encrypt file (if needed) and upload. No archive
					for _, f := range fileList {

						if f, err = encryptAndUpload(f, backupConfig.EncryptionKey); err != nil {
							glg.Error(err)
							return
						}

						if err = os.Remove(f); err != nil {
							glg.Error(err)
							return
						}

						if !runFlag.IsSet() {
							glg.Warn("Run flag disabled")
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
