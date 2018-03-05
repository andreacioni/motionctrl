package notify

import (
	"fmt"
	"sync"

	"github.com/andreacioni/motionctrl/utils"

	"github.com/marusama/semaphore"

	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/config"
)

type NotifyService interface {
	Authenticate() error
	Notify(string) error
	Stop() error
}

type State string

const (
	TelegramNotifyMethod = "telegram"
	TestMockMethod       = "mock"
)

var (
	notifyConfiguration config.Notify

	notifyService NotifyService

	photoLimitSemaphore semaphore.Semaphore

	mu sync.Mutex
)

func Init(conf config.Notify) {
	mu.Lock()
	defer mu.Unlock()

	if notifyService == nil {
		var err error
		notifyConfiguration = conf

		glg.Debugf("Initializing notify service: %+v", conf)

		notifyService, err = buildNotifyService(conf)

		if err != nil {
			glg.Warn("Notify service won't be active")
		} else {
			if err = notifyService.Authenticate(); err != nil {
				glg.Errorf("Cannot authenticate to '%s' service: %v", conf.Method, err)
			} else {
				photoLimitSemaphore = semaphore.New(conf.Photo)
			}
		}
	}
}

func Shutdown() {
	mu.Lock()
	defer mu.Unlock()

	if notifyService != nil {
		notifyService.Stop()
		notifyService = nil
	}

	notifyConfiguration = config.Notify{}
}

func MotionDetectedStart() {
	mu.Lock()
	defer mu.Unlock()

	if notifyService != nil {
		photoLimitSemaphore.Release(notifyConfiguration.Photo)
	} else {
		glg.Warn("No notify service is available")
	}
}

func MotionDetectedStop() {
	mu.Lock()
	defer mu.Unlock()

	if notifyService != nil {
		for photoLimitSemaphore.GetCount() > 0 && !photoLimitSemaphore.TryAcquire(1) { //Clear all
		}
	} else {
		glg.Warn("No notify service is available")
	}
}

func PhotoSaved(filepath string) {
	mu.Lock()
	defer mu.Unlock()

	if notifyService != nil {
		if photoLimitSemaphore.TryAcquire(1) {
			if err := notifyService.Notify(filepath); err != nil {
				glg.Errorf("Failed to send notify: %v", err)
			}
		} else {
			glg.Warnf("Photo limit reached, this one won't be sent")
		}
	} else {
		glg.Warn("No notify service is available")
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
