package motion

import (
	"fmt"
	"path/filepath"

	"github.com/andreacioni/motionctrl/utils"
)

func Snapshot() (string, error) {
	snapFile, err := webControlGet("/action/snapshot", func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(SnapshotDetectionRegex, body) {
			return "", fmt.Errorf("unable to take snapshot (%s)", body)
		}
		return filepath.Join(ConfigGetRO(ConfigTargetDir), "lastsnap.jpg"), nil //TODO jpg not the only extension
	})

	return snapFile.(string), err
}
