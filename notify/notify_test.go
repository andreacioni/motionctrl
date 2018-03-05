package notify

import (
	"testing"

	"github.com/andreacioni/motionctrl/config"
)

func TestNotify(t *testing.T) {
	Init(config.Notify{Method: "mock"})
	Shutdown()
}

func TestNotify(t *testing.T) {
	Init(config.Notify{Method: "mock"})
	Shutdown()
}
