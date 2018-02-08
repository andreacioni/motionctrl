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

	err := CheckInstall()

	if err != nil {
		glg.Fatalf("Motion not found (%s)", err)
	}

	if config.Get().MotionConfigFile != "" {
		glg.Info("Motion config file specified", config.Get().MotionConfigFile)
	} else {
		glg.Fatalf("Motion config file is not defined in configuration, %s can't start without it", version.Name)
	}

	glg.Infof("Loading motion configuration from %s...", config.Get().MotionConfigFile)

	err = Load(config.Get().MotionConfigFile)

	if err != nil {
		glg.Fatalf("Failed to load motion configuration file (%s)", err)
	}
}

//CheckInstall will check if motion is available and ready to be controlled. If motion isn't available the program will exit showing an error
func CheckInstall() error {
	err := exec.Command("motion", "-h").Run()

	//TODO unfortunatelly motion doesn't return 0 when invoked with the "-h" parameter
	if err != nil || err.Error() != "exit status 1" {
		return err
	}

	return nil

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
