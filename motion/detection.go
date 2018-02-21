package motion

import (
	"fmt"

	"github.com/andreacioni/motionctrl/utils"
)

const (
	DetectionStatusRegex  = "Camera [0-9]+ Detection status (ACTIVE|PAUSE)"
	DetectionResumedRegex = "Camera [0-9]+ Detection resumed\nDone\n"
	DetectionPausedRegex  = "Camera [0-9]+ Detection paused\nDone\n"
	DetectionActiveRegex  = "Camera [0-9]+ Detection status ACTIVE"
	DetectionPauseRegex   = "Camera [0-9]+ Detection status PAUSE"
)

func IsMotionDetectionEnabled() (bool, error) {
	ret, err := webControlGet("/detection/status", func(body string) (interface{}, error) {
		status := utils.RegexFirstSubmatchString(DetectionStatusRegex, body)
		if status == "ACTIVE" {
			return true, nil
		} else if status == "PAUSE" {
			return false, nil
		} else {
			return false, fmt.Errorf("unknown status string: %s", status)
		}
	})

	if err != nil {
		return false, err
	}

	return ret.(bool), err
}

func EnableMotionDetection() error {

	_, err := webControlGet("/detection/start", func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(DetectionResumedRegex, body) {
			return nil, fmt.Errorf("unable to enable motion detection (%s)", body)
		}
		return nil, nil
	})

	return err
}

func DisableMotionDetection() error {

	_, err := webControlGet("/detection/pause", func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(DetectionPausedRegex, body) {
			return nil, fmt.Errorf("unable to disable motion detection (%s)", body)
		}
		return nil, nil
	})

	return err
}
