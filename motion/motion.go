package motion

import (
	"fmt"
	"net/http"
	"os/exec"
	"sync"

	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/version"
	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

var (
	sMutex sync.Mutex

	motionConfigFile string
)

func Init(configFile string, autostart bool, detection bool) error {
	sMutex.Lock()
	defer sMutex.Unlock()

	if err := checkInstall(); err != nil {
		return fmt.Errorf("Motion not found: %v", err)
	}

	if started, err := checkStarted(); err == nil {
		if started {
			glg.Warn("Motion started before %s", version.Name)
		}
	} else {
		return fmt.Errorf("Unable to check is motion is running: %v", err)
	}

	if err := loadConfig(configFile); err != nil {
		return fmt.Errorf("Failed to load motion configuration: %v", err)
	}

	motionConfigFile = configFile

	if autostart {
		glg.Infof("Starting motion")

		if err := startMotion(detection); err != nil {
			return fmt.Errorf("Unable to start motion: %v", err)
		}
	}

	return nil
}

func GetStreamBaseURL() string {
	return fmt.Sprintf("http://%s:%s", config.BaseAddress, readOnlyConfig[ConfigStreamPort])
}

func GetBaseURL() string {
	return fmt.Sprintf("http://%s:%s", config.BaseAddress, readOnlyConfig[ConfigWebControlPort])
}

//CheckInstall will check if motion is available and ready to be controlled. If motion isn't available the program will exit showing an error
func checkInstall() error {
	err := exec.Command("motion", "-h").Run()

	//TODO unfortunatelly motion doesn't return 0 when invoked with the "-h" parameter
	if err != nil && err.Error() != "exit status 1" {
		return err
	}

	return nil

}

func webControlGet(path string, callback func(string) (interface{}, error)) (interface{}, error) {
	var err error
	var ret interface{}

	resp, body, errs := gorequest.New().Get(GetBaseURL() + "/0" + path).End()

	if errs == nil {
		glg.Debugf("Response body: %s", body)
		if resp.StatusCode == http.StatusOK {
			ret, err = callback(body)
		} else {
			ret, err = nil, fmt.Errorf("request failed with code: %d", resp.StatusCode)
		}
	} else {
		ret, err = nil, errs[0] //TODO errs[0] not the best
	}

	return ret, err
}
