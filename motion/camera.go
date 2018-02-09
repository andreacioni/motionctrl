package motion

import (
	"fmt"

	"../config"

	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

func IsMotionDetectionEnabled() (bool, error) {
	glg.Debugf("Connecting to: %s", fmt.Sprintf("%s:%s", config.BaseAddress, motionConfMap[WebControlPort]))

	resp, body, errs := gorequest.New().Get(fmt.Sprintf("%s:%s/0/detection/status", config.BaseAddress, motionConfMap[WebControlPort])).End()

	if errs == nil {
		if resp.StatusCode == 200 {

		} else {
			return false, fmt.Errorf("request failed")
		}
	} else {
		return false, errs
	}
}
