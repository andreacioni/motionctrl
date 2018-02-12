package motion

import (
	"fmt"
	"sync"

	"../config"
)

const (
	DetectionStatusRegex = "Camera [0-9]+ Detection status (ACTIVE|PAUSE)"
	DoneRegex            = "\nDone"
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
