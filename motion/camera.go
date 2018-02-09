package motion

import (
	"fmt"

	"../config"
	"../utils"

	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

const (
	DETECTION_STATUS_REGEX = "Camera [0-9]+ Detection status (ACTIVE|PAUSE)"
)

func IsMotionDetectionEnabled() (bool, error) {
	glg.Debugf("Connecting to: %s", fmt.Sprintf("%s:%s", config.BaseAddress, motionConfMap[WebControlPort]))

	resp, body, errs := gorequest.New().Get(fmt.Sprintf("%s:%s/0/detection/status", config.BaseAddress, motionConfMap[WebControlPort])).End()

	if errs == nil {
		if resp.StatusCode == 200 {

			if utils.RegexFirstSubmatchString(DETECTION_STATUS_REGEX, body) == "ACTIVE" {

			} else if utils.RegexFirstSubmatchString(DETECTION_STATUS_REGEX, body) == "PAUSE" {

			} else {

			}

		} else {
			return false, fmt.Errorf("request failed")
		}
	} else {
		return false, errs
	}
}
