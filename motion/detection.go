package motion

import (
	"fmt"
	"net/http"

	"../config"
	"../utils"
	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

func IsMotionDetectionEnabled() (bool, error) {
	var err error
	ret := false

	resp, body, errs := gorequest.New().Get(GetBaseURL() + "/detection/status").End()

	if errs == nil {
		if resp.StatusCode == http.StatusOK {
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

	command := ""
	if enable {
		command = fmt.Sprintf(GetBaseURL()+"/detection/start", config.BaseAddress, motionConfMap[WebControlPort])
	} else {
		command = fmt.Sprintf(GetBaseURL()+"/detection/pause", config.BaseAddress, motionConfMap[WebControlPort])
	}

	resp, body, errs := gorequest.New().Get(command).End()

	if errs == nil {
		if resp.StatusCode == http.StatusOK {
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
