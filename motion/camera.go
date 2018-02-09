package motion

import (
	"fmt"

	"../config"
	"../utils"

	"github.com/parnurzeal/gorequest"
)

const (
	DetectionStatusRegex = "Camera [0-9]+ Detection status (ACTIVE|PAUSE)"
	DoneRegex            = "\nDone"
)

func IsMotionDetectionEnabled() (bool, error) {

	resp, body, errs := gorequest.New().Get(fmt.Sprintf("http://%s:%s/0/detection/status", config.BaseAddress, motionConfMap[WebControlPort])).End()

	if errs == nil {
		if resp.StatusCode == 200 {
			status := utils.RegexFirstSubmatchString(DetectionStatusRegex, body)
			if status == "ACTIVE" {
				return true, nil
			} else if status == "PAUSE" {
				return false, nil
			} else {
				return false, fmt.Errorf("unknown status string: %s", status)
			}

		} else {
			return false, fmt.Errorf("request failed with code: %d", resp.StatusCode)
		}
	} else {
		return false, errs[0] //TODO errs[0] not the best
	}
}

func EnableMotionDetection(enable bool) error { //TODO check for motion is running
	command := ""
	if enable {
		command = fmt.Sprintf("http://%s:%s/0/detection/start", config.BaseAddress, motionConfMap[WebControlPort])
	} else {
		command = fmt.Sprintf("http://%s:%s/0/detection/pause", config.BaseAddress, motionConfMap[WebControlPort])
	}

	resp, body, errs := gorequest.New().Get(command).End()

	if errs == nil {
		if resp.StatusCode == 200 {

			if utils.RegexMustMatch(DoneRegex, body) {
				return nil
			}
			return fmt.Errorf("unable to enable motion detection (%s)", body)
		}
		return fmt.Errorf("request failed with code: %d", resp.StatusCode)
	}

	return errs[0] //TODO errs[0] not the best
}
