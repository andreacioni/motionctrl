package motion

import (
	"fmt"

	"../config"
	"../utils"

	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

const (
	DetectionStatusRegex = "Camera [0-9]+ Detection status (ACTIVE|PAUSE)"
	DoneRegex            = "\nDone"
)

func Shutdown() error {
	var err error
	mu.Lock()
	defer mu.Unlock()

	glg.Debug("Stoping motion")

	if started {
		resp, body, errs := gorequest.New().Get(getBaseURL() + "/0/action/quit").End()

		if errs == nil {
			if resp.StatusCode == 200 {
				glg.Debugf("Response body: %s", body)

				if utils.RegexMustMatch(DoneRegex, body) {
					err = waitDie()

					if err == nil {
						started = false
					}
				} else {
					err = fmt.Errorf("failed to quit motion: %s", body)
				}

			} else {
				err = fmt.Errorf("request failed with code: %d", resp.StatusCode)
			}

		} else {
			err = errs[0]
		}
	} else {
		glg.Warn("motion is already stopped")
	}

	return err
}

func Restart() error {
	var err error
	mu.Lock()
	defer mu.Unlock()

	glg.Debug("Restarting motion")

	if started {
		resp, body, errs := gorequest.New().Get(getBaseURL() + "/0/action/restart").End()

		if errs == nil {
			if resp.StatusCode == 200 {
				glg.Debugf("Response body: %s", body)

				if utils.RegexMustMatch(DoneRegex, body) {
					err = waitDie()

					if err == nil {
						err = waitLive()
					}

				} else {
					err = fmt.Errorf("failed to restart motion: %s", body)
				}
			} else {
				err = fmt.Errorf("request failed with code: %d", resp.StatusCode)
			}

		} else {
			err = errs[0]
		}
	} else {
		glg.Warn("motion is not started")
	}

	return err
}

func IsMotionDetectionEnabled() (bool, error) {
	ret, err := false, error(nil)
	mu.Lock()
	defer mu.Unlock()

	resp, body, errs := gorequest.New().Get(getBaseURL() + "0/detection/status").End()

	if errs == nil {
		if resp.StatusCode == 200 {
			glg.Debugf("Response body: %s", body)

			status := utils.RegexFirstSubmatchString(DetectionStatusRegex, body)
			if status == "ACTIVE" {
				ret, err = true, nil
			} else if status == "PAUSE" {
				ret, err = false, nil
			} else {
				ret, err = false, fmt.Errorf("unknown status string: %s", status)
			}

		} else {
			ret, err = false, fmt.Errorf("request failed with code: %d", resp.StatusCode)
		}
	} else {
		ret, err = false, errs[0] //TODO errs[0] not the best
	}

	return ret, err
}

func EnableMotionDetection(enable bool) error { //TODO check for motion is running
	err := error(nil)
	mu.Lock()
	defer mu.Unlock()

	command := ""
	if enable {
		command = fmt.Sprintf(getBaseURL()+"/0/detection/start", config.BaseAddress, motionConfMap[WebControlPort])
	} else {
		command = fmt.Sprintf(getBaseURL()+"/0/detection/pause", config.BaseAddress, motionConfMap[WebControlPort])
	}

	resp, body, errs := gorequest.New().Get(command).End()

	if errs == nil {
		if resp.StatusCode == 200 {
			glg.Debugf("Response body: %s", body)

			if utils.RegexMustMatch(DoneRegex, body) {
				err = fmt.Errorf("unable to enable motion detection (%s)", body)
			}
		} else {
			err = fmt.Errorf("request failed with code: %d", resp.StatusCode)
		}

	}

	return err //TODO errs[0] not the best
}
