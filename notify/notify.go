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

	var err error
	notifyConfiguration = conf

	glg.Debugf("Initializing notify service: %+v", notifyConfiguration)

	notifyService, err = buildNotifyService(notifyConfiguration)

	if err != nil {
		glg.Warn("Notify service won't be active")
	} else {
		if err = notifyService.Authenticate(); err != nil {
			glg.Errorf("Cannot authenticate to '%s' service: %v", conf.Method, err)
		} else {
			photoLimitSemaphore = semaphore.New(0)
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
	}
}

func MotionDetectedStop() {
	mu.Lock()
	defer mu.Unlock()

	if notifyService != nil {
		for !photoLimitSemaphore.TryAcquire(1) { //Clear all
		}
	}
}

func PhotoSaved(filepath string) {
	mu.Lock()

	if notifyService != nil {
		mu.Unlock()
		if photoLimitSemaphore.TryAcquire(1) {
			if err := notifyService.Notify(filepath); err != nil {
				glg.Errorf("Failed to send notify: %v", err)
			}
		} else {
			glg.Warnf("Photo limit reached, this one won't be sent")
		}
	} else {
		mu.Unlock()
	}
}

func buildNotifyService(conf config.Notify) (NotifyService, error) {
	switch conf.Method {
	case TelegramNotifyMethod:
		if chatids, err := utils.ToInt64Slice(conf.To); err != nil {
			return nil, fmt.Errorf("Failed to convert chat IDs from string to int: %v", err)
		} else {
			return &TelegramNotifyService{apiToken: conf.Token, chatIds: chatids}, nil
		}
	case TestMockMethod:
		return &MockNotifyService{}, nil
	default:
		return nil, fmt.Errorf("Notify method not found or invalid")
	}
}
