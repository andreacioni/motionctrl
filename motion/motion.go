package motion

import (
	"fmt"
	"net/http"
	"sync"

	"../config"
	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"
)

var (
	mu               sync.Mutex
	started          bool
	motionConfigFile string
)

func GetStreamBaseURL() string {
	return fmt.Sprintf("http://%s:%s", config.BaseAddress, motionConfMap[StreamPort])
}

func GetBaseURL() string {
	return fmt.Sprintf("http://%s:%s/0", config.BaseAddress, motionConfMap[WebControlPort])
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
