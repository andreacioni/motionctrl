package motion

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/version"
	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

var (
	mu               sync.Mutex
	started          bool
	motionConfigFile string
)

func Init(configFile string, autostart bool, detection bool) error {

	err := checkInstall()

	if err != nil {
		return fmt.Errorf("Motion not found (%s)", err)
	}

	if configFile != "" {
		_, err = os.Stat(configFile)
		if err != nil {
			return fmt.Errorf("Cannot open file %s", configFile)
		}

		glg.Infof("Motion config file specified: %s", configFile)
	} else {
		return fmt.Errorf("Motion config file is not defined in configuration, %s can't start without it", version.Name)
	}

	glg.Infof("Loading motion configuration from %s...", configFile)

	err = loadConfig(configFile)

	if err != nil {
		return fmt.Errorf("Failed to load motion configuration file (%s)", err)
	}

	_, err = readPid()
	if err == nil {
		return fmt.Errorf("Motion is started before %s. Kill motion and retry", version.Name)
	}

	motionConfigFile = configFile

	if autostart {
		glg.Infof("Starting motion")
		err = Startup(detection)

		if err != nil {
			return fmt.Errorf("Unable to startup motion (%s)", err)
		}
	}

	return nil
}

func GetStreamBaseURL() string {
	return fmt.Sprintf("http://%s:%s", config.BaseAddress, readOnlyConfig[ConfigStreamPort])
}

func GetBaseURL() string {
	return fmt.Sprintf("http://%s:%s/0", config.BaseAddress, readOnlyConfig[ConfigWebControlPort])
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

	resp, body, errs := gorequest.New().Get(GetBaseURL() + path).End()

	if errs == nil {
		if resp.StatusCode == http.StatusOK {
			glg.Debugf("Response body: %s", body)
			ret, err = callback(body)
		} else {
			ret, err = nil, fmt.Errorf("request failed with code: %d", resp.StatusCode)
		}
	} else {
		ret, err = nil, errs[0] //TODO errs[0] not the best
	}

	return ret, err
}
