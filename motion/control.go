package motion

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

func Startup(motionDetectionStartup bool) error {
	var err error
	mu.Lock()
	defer mu.Unlock()

	glg.Debugf("Starting motion (detection enabled: %t)", motionDetectionStartup)

	if !started {
		err = startMotion(motionDetectionStartup)
		if err == nil {
			started = true
		}
	} else {
		glg.Warn("motion is already started")
	}

	return err
}

func Shutdown() error {
	var err error
	mu.Lock()
	defer mu.Unlock()

	glg.Debug("Stopping motion")

	if started {
		stopMotion()
	} else {
		glg.Warn("motion is already stopped")
	}

	started = false

	return err
}

func Restart() error {
	var err error
	var detection bool
	mu.Lock()
	defer mu.Unlock()

	glg.Debug("Restarting motion")

	if started {
		detection, err = IsMotionDetectionEnabled()
		if err == nil {
			err = stopMotion()
			if err == nil {
				err = startMotion(detection)
			}
		}

	} else {
		err = fmt.Errorf("motion is not running")
	}

	return err
}

func IsStarted() bool {
	mu.Lock()
	defer mu.Unlock()
	return started
}

func readPid() (int, error) {
	raw, err := ioutil.ReadFile(readOnlyConfig[ConfigProcessIdFile])

	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(raw[:len(raw)-1]))

	if err != nil {
		return 0, err
	}

	return pid, err
}

func startMotion(motionDetectionStartup bool) error {
	var err error
	if motionDetectionStartup {
		err = exec.Command("motion", "-b", "-c", motionConfigFile).Run()
	} else {
		err = exec.Command("motion", "-b", "-m", "-c", motionConfigFile).Run()
	}

	if err == nil {
		err = waitLive()
	}

	return err
}

func stopMotion() error {
	pid, err := readPid()

	if err == nil {
		glg.Debugf("Going to kill motion (PID: %d)", pid)
		err = exec.Command("kill", "-2", fmt.Sprint(pid)).Run()
		if err == nil {
			err = waitDie()
		}
	}

	return err
}

func waitDie() error {
	i, secs := 0, 15
	for _, err := os.Stat(readOnlyConfig[ConfigProcessIdFile]); err == nil && i < secs; _, err = os.Stat(readOnlyConfig[ConfigProcessIdFile]) {
		glg.Debugf("Waiting motion exits (attempts: %d/%d)", i, secs)
		time.Sleep(time.Second)
		i++
	}

	if i == secs {
		return fmt.Errorf("motion is alive after %d seconds", secs)
	}

	return nil
}

func waitLive() error {
	req := gorequest.New()
	i, secs := 0, 15
	for _, _, errs := req.Get(GetBaseURL()).End(); errs != nil && i < secs; _, _, errs = req.Get(GetBaseURL()).End() {
		glg.Debugf("Waiting motion to become available (attempts: %d/%d)", i, secs)
		time.Sleep(time.Second)
		i++
	}

	if i == secs {
		return fmt.Errorf("motion is not ready after %d seconds", secs)
	}

	return nil
}
