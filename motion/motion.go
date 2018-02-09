package motion

import (
	"fmt"
	"os"
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
		_, err = os.Stat(config.Get().MotionConfigFile)
		if err != nil {
			glg.Info("Motion config file specified", config.Get().MotionConfigFile)
		} else {
			glg.Fatalf("Cannot open file %s", config.Get().MotionConfigFile)
		}

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
	if err != nil && err.Error() != "exit status 1" {
		return err
	}

	return nil

}

func Startup(motionDetectionStartup bool) error {
	mu.Lock()
	defer mu.Unlock()

	var err error

	if motionDetectionStartup {
		err = exec.Command("motion", "-b", "-c", config.Get().MotionConfigFile).Run()
	} else {
		err = exec.Command("motion", "-b", "-m", "-c", config.Get().MotionConfigFile).Run()
	}

	if err != nil {
		return err
	}

	started = true

	return nil
}

func Shutdown() error {
	mu.Lock()
	defer mu.Unlock()

	if started {
		err := exec.Command("kill", "-2", fmt.Sprintf("$(cat %s)", motionConfMap[ProcessIdFile]), config.Get().MotionConfigFile).Run()
		if err != nil {
			return err
		}
	}

	started = false

	return nil
}

func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}
