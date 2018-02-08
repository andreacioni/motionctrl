package motion

import (
	"os/exec"
	"sync"

	"../config"
	"../version"
	"github.com/kpango/glg"
)

var (
	mu      sync.Mutex
	started bool
)

func Init() {
	if config.Conf.MotionConfig != "" {
		glg.Infof("Motion config file specified", config.Conf.MotionConfig)
	} else {
		glg.Fatalf("Motion config file is not defined in configuration, %s can't start without it", version.Name)
	}
}

//CheckInstall will check if motion is available and ready to be controlled. If motion isn't available the program will exit showing an error
func CheckInstall() {
	err := exec.Command("motion", "-h").Run()

	//TODO unfortunatelly motion doesn't return 0 when invoked with the "-h" parameter
	if err != nil && err.Error() != "exit status 1" {
		glg.Fatalf("Motion not found (%s)", err)
	}

	glg.Debug("Motion found")
}

func Startup(motionDetectionStartup bool) {
	mu.Lock()
	defer mu.Unlock()
	started = true
}

func Shutdown() {
	mu.Lock()
	defer mu.Unlock()
	started = false
}

func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}
