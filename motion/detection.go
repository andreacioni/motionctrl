package motion

import (
	"fmt"

	"../config"
	"../utils"
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

	if ret != nil {
		return false, err
	}

	return ret.(bool), err
}

func EnableMotionDetection(enable bool) error { //TODO check for motion is running
	path := ""
	if enable {
		path = fmt.Sprintf(GetBaseURL()+"/detection/start", config.BaseAddress, motionConfMap[WebControlPort])
	} else {
		path = fmt.Sprintf(GetBaseURL()+"/detection/pause", config.BaseAddress, motionConfMap[WebControlPort])
	}

	_, err := webControlGet(path, func(body string) (interface{}, error) {
		if utils.RegexMustMatch(DoneRegex, body) {
			return nil, fmt.Errorf("unable to enable motion detection (%s)", body)
		}
		return nil, nil
	})

	return err
}
