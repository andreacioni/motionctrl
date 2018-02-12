package motion

import (
	"fmt"

	"../config"
	"../utils"
)

func IsMotionDetectionEnabled() (bool, error) {
	ret, err := WebControlGet("/detection/status", func(body string) (interface{}, error) {
		status := utils.RegexFirstSubmatchString(DetectionStatusRegex, body)
		if status == "ACTIVE" {
			return true, nil
		} else if status == "PAUSE" {
			return false, nil
		} else {
			return false, fmt.Errorf("unknown status string: %s", status)
		}
	})

	return ret.(bool), err
}

func EnableMotionDetection(enable bool) error { //TODO check for motion is running
	path := ""
	if enable {
		path = fmt.Sprintf(GetBaseURL()+"/detection/start", config.BaseAddress, motionConfMap[WebControlPort])
	} else {
		path = fmt.Sprintf(GetBaseURL()+"/detection/pause", config.BaseAddress, motionConfMap[WebControlPort])
	}

	return WebControlGet(path, func(body string) (interface{}, error) {
		var err error
		if utils.RegexMustMatch(DoneRegex, body) {
			err = fmt.Errorf("unable to enable motion detection (%s)", body)
		}
		return nil, err
	})
}
