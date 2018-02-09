package motion

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"../version"
	"github.com/kpango/glg"
)

var (
	mu               sync.Mutex
	started          bool
	motionConfigFile string
)

func Init(configFile string) {

	err := CheckInstall()

	if err != nil {
		glg.Fatalf("Motion not found (%s)", err)
	}

	if configFile != "" {
		_, err = os.Stat(configFile)
		if err != nil {
			glg.Fatalf("Cannot open file %s", configFile)
		} else {

			glg.Infof("Motion config file specified: %s", configFile)
		}

	} else {
		glg.Fatalf("Motion config file is not defined in configuration, %s can't start without it", version.Name)
	}

	glg.Infof("Loading motion configuration from %s...", configFile)

	err = Load(configFile)

	if err != nil {
		glg.Fatalf("Failed to load motion configuration file (%s)", err)
	}

	_, err = readPid()
	if err == nil {
		glg.Fatalf("Motion is started before %s. Kill motion and retry", version.Name)
	}

	motionConfigFile = configFile
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
	if !started {
		var err error

		if motionDetectionStartup {
			err = exec.Command("motion", "-b", "-c", motionConfigFile).Run()
		} else {
			err = exec.Command("motion", "-b", "-m", "-c", motionConfigFile).Run()
		}

		if err != nil {
			return err
		}
	} else {
		glg.Warn("motion is already started")
	}

	started = true

	return nil
}

func Shutdown() error {
	mu.Lock()
	defer mu.Unlock()

	if started {

		pid, err := readPid()

		if err == nil {
			glg.Debugf("Going to kill motion (PID: %d)", pid)
			err := exec.Command("kill", "-2", fmt.Sprint(pid)).Run()
			if err != nil {
				return fmt.Errorf("failed to kill motion instance: %s", err)
			}
		} else {
			return fmt.Errorf("cannot read pid of motion process: %s", err)
		}

	} else {
		glg.Warn("motion is already stopped")
	}

	started = false

	return nil
}

func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}

func readPid() (int, error) {
	raw, err := ioutil.ReadFile(motionConfMap[ProcessIdFile])

	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(raw[:len(raw)-1]))

	if err != nil {
		return 0, err
	}

	return pid, err
}
