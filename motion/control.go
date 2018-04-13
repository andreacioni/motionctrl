package motion

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/andreacioni/motionctrl/utils"

	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

const (
	waitLiveRegex = "Motion"
)

var (
	sMutex sync.Mutex
)

func Startup(motionDetectionStartup bool) error {
	sMutex.Lock()
	defer sMutex.Unlock()

	glg.Debugf("Starting motion (detection enabled: %t)", motionDetectionStartup)

	var err error
	var started bool

	if started, err = checkStarted(); err == nil {
		if !started {
			if err = startMotion(motionDetectionStartup); err != nil {
				return err
			}
		} else {
			glg.Warn("motion is already started")
		}
	} else {
		err = fmt.Errorf("unable to check if motion is started: %v", err)
	}

	return err
}

func Shutdown() error {
	sMutex.Lock()
	defer sMutex.Unlock()

	var err error
	var started bool

	glg.Debug("Stopping motion")

	if started, err = checkStarted(); err == nil {
		if started {
			err = stopMotion()
		} else {
			glg.Warn("motion is already stopped")
		}
	} else {
		err = fmt.Errorf("unable to check if motion is started: %v", err)
	}

	return err
}

func Restart() error { //TODO started consistency isn't guaranteed
	sMutex.Lock()
	defer sMutex.Unlock()

	var err error
	var detection, started bool

	glg.Debug("Restarting motion")

	if started, err = checkStarted(); err == nil {
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
	} else {
		err = fmt.Errorf("unable to check if motion is started: %v", err)
	}

	return err
}

func IsStarted() (bool, error) {
	sMutex.Lock()
	defer sMutex.Unlock()
	return checkStarted()
}

func checkStarted() (bool, error) {
	pid, err := readPid()

	if err != nil {
		return false, nil
	}

	if check := isRunningPID(pid); !check {
		/*
			Motion pid file isn't valid, probably it belongs to a Motion instance that
			wasn't terminate correctly. Removing old file.
		*/

		glg.Warn("Removing old PID file")

		if err := os.Remove(readOnlyConfig[ConfigProcessIdFile]); err != nil {
			return false, fmt.Errorf("Failed to remove invalid PID file: %v", err)
		}

		return false, nil
	}

	return true, nil
}

func isRunningPID(pid int) bool {
	/*
		from docs: On Unix systems, FindProcess always succeeds and returns
		a Process for the given pid, regardless of whether the process exists.
	*/
	proc, _ := os.FindProcess(pid)

	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return false //err == syscall.EPERM
	}

	return true
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
	if pid, err := readPid(); err == nil {
		glg.Debugf("Going to kill motion (PID: %d)", pid)
		/*
			from docs: On Unix systems, FindProcess always succeeds and returns
			a Process for the given pid, regardless of whether the process exists.
		*/
		if proc, err := os.FindProcess(pid); err == nil {
			if err := proc.Signal(syscall.SIGINT); err != nil {
				return err
			}
		} else {
			return err
		}

	} else {
		return err
	}

	return waitDie()
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
	var err error
	var pid int

	if pid, err = readPid(); err != nil {
		return err
	}
	glg.Debugf("motion started (PID: %d", pid)

	req := gorequest.New()
	i, secs := 0, 15
	for resp, body, errs := req.Get(GetBaseURL()).End(); i < secs; resp, body, errs = req.Get(GetBaseURL()).End() {
		if errs == nil && resp.StatusCode == http.StatusOK && utils.RegexMustMatch(waitLiveRegex, body) {
			break
		}

		glg.Debugf("Waiting motion to become available (attempts: %d/%d)", i, secs)
		time.Sleep(time.Second)
		i++
	}

	if i == secs {
		return fmt.Errorf("motion is not ready after %d seconds", secs)
	}

	return nil
}
