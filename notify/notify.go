package notify

import (
	"fmt"
	"sync"

	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/config"
)

type NotifyService interface {
	Authenticate() error
	Notify(string) error
}

type State string

const (
	TelegramNotifyMethod = "telegram"
	TestMockMethod       = "mock"
)

var (
	notifyConfiguration config.Notify

	notifyService NotifyService

	backupStateMutex sync.Mutex
)

func Init(conf config.Notify) {
	var err error
	notifyConfiguration = conf

	glg.Debugf("Initializing notify service: %+v", notifyConfiguration)

	notifyService, err = buildNotifyService(notifyConfiguration)

	if err != nil {
		glg.Warn("Notify service won't be active")
	} else {
		if err = notifyService.Authenticate(); err != nil {
			glg.Errorf("Cannot authenticate to '%s' service: %v", conf.Method, err)
		}
	}
}

func buildNotifyService(conf config.Notify) (NotifyService, error) {
	switch conf.Method {
	case TelegramNotifyMethod:
		return &TelegramNotifyService{apiToken: conf.Token, chatIds: conf.To}, nil
	case TestMockMethod:
		return &MockNotifyService{}, nil
	default:
		return nil, fmt.Errorf("Notify method not found or invalid")
	}
}
