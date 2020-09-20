package notify

import (
	"fmt"
	"sync"
	"time"

	"github.com/andreacioni/motionctrl/utils"

	"github.com/abiosoft/semaphore"

	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/config"
)

type NotifyService interface {
	Authenticate() error
	Notify(string, string) error
	Stop() error
}

type State string

const (
	TelegramNotifyMethod = "telegram"
	TestMockMethod       = "mock"
)

var (
	notifyConfiguration config.Notify

	nMutex        sync.Mutex
	notifyService NotifyService
	active        bool

	photoLimitSemaphore *semaphore.Semaphore
)

func Init(conf config.Notify) error {
	nMutex.Lock()
	defer nMutex.Unlock()

	var err error

	if notifyService == nil {
		notifyConfiguration = conf

		glg.Debugf("Initializing notify service: %+v", conf)

		notifyService, err = buildNotifyService(conf)

		if err != nil {
			glg.Warnf("Notify service won't be active: %v", err)
		} else {
			if err = notifyService.Authenticate(); err != nil {
				err = fmt.Errorf("Cannot authenticate to '%s' service: %v", conf.Method, err)
				notifyService = nil
			} else {

				active = true

				if conf.Photo > 0 {
					photoLimitSemaphore = semaphore.New(notifyConfiguration.Photo)
					photoLimitSemaphore.DrainPermits()
				} else {
					glg.Warn("No photo will be sent when motion is detected")
				}
			}
		}
	} else {
		err = fmt.Errorf("Notify service is already running")
	}

	return err
}

func Shutdown() {
	nMutex.Lock()
	defer nMutex.Unlock()

	glg.Info("Shuting down notify service")

	if notifyService != nil {
		notifyService.Stop()
		notifyService = nil
	}

	drainPhotoSemaphore()
	photoLimitSemaphore = nil

	notifyConfiguration = config.Notify{}

	active = false
}

func MotionDetectedStart() {
	nMutex.Lock()
	defer nMutex.Unlock()

	if notifyService != nil {
		if active {
			resetPhotoSemaphore()
			notifyService.Notify(notifyConfiguration.Message, "")
		} else {
			glg.Warn("Notify service deactivated")
		}

	} else {
		glg.Warn("No notify service is available")
	}
}

func MotionDetectedStop() {
	nMutex.Lock()
	defer nMutex.Unlock()

	if notifyService != nil {
		if active {
			drainPhotoSemaphore()
		} else {
			glg.Warn("Notify service deactivated")
		}
	} else {
		glg.Warn("No notify service is available")
	}
}

func PhotoSaved(filepath string) {
	nMutex.Lock()
	defer nMutex.Unlock()

	if notifyService != nil {
		if active {
			if photoLimitSemaphore != nil {
				if photoLimitSemaphore.AcquireWithin(1, 10*time.Microsecond) { //TODO improve this
					if err := notifyService.Notify("", filepath); err != nil {
						glg.Errorf("Failed to send notify: %v", err)
					} else {
						glg.Debugf("Sent notify (image; %s)", filepath)
					}
				} else {
					glg.Warnf("Photo limit reached, this one won't be sent")
				}
			} else {
				glg.Warn("Unable to send picture, 'photo' parameter not defined?")
			}

		} else {
			glg.Warn("Notify service deactivated")
		}

	} else {
		glg.Warn("No notify service is available")
	}
}

func IsReady() bool {
	nMutex.Lock()
	defer nMutex.Unlock()

	return notifyService != nil
}

func IsActive() bool {
	nMutex.Lock()
	defer nMutex.Unlock()

	return active
}

func SetActive(status bool) error {
	nMutex.Lock()
	defer nMutex.Unlock()

	if notifyService == nil {
		return fmt.Errorf("No notify service is currently configured")
	}

	glg.Debugf("Setting notify service active: %t", status)

	active = status

	return nil
}

func resetPhotoSemaphore() {
	if photoLimitSemaphore != nil {
		photoLimitSemaphore.ReleaseMany(notifyConfiguration.Photo)
	}
}

func drainPhotoSemaphore() {
	if photoLimitSemaphore != nil {
		photoLimitSemaphore.DrainPermits()
	}
}

func buildNotifyService(conf config.Notify) (NotifyService, error) {
	switch conf.Method {
	case TelegramNotifyMethod: //TODO following statemets should be moved to an appropriate factory
		var chatids []int64
		var err error
		if chatids, err = utils.ToInt64Slice(conf.To); err != nil {
			return nil, fmt.Errorf("Failed to convert chat IDs from string to int: %v", err)
		}
		return &TelegramNotifyService{apiToken: conf.Token, chatIds: chatids}, nil
	case TestMockMethod:
		return &MockNotifyService{}, nil
	default:
		return nil, fmt.Errorf("Notify method not found or invalid")
	}
}
