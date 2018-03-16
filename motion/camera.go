package motion

import (
	"fmt"
	"path/filepath"

	"github.com/andreacioni/motionctrl/utils"
)

func Snapshot() (string, error) {

	snapExt, err := ConfigGet(ConfigPictureType)

	//TODO snapshot file not ready when function return. This cause some 404
	snapFile, err := webControlGet("/action/snapshot", func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(SnapshotDetectionRegex, body) {
			return "", fmt.Errorf("unable to take snapshot (%s)", body)
		}

		switch snapExt[ConfigPictureType] {
		case "ppm":
			return filepath.Join(ConfigGetRO(ConfigTargetDir), "lastsnap.ppm"), nil
		case "webp":
			return filepath.Join(ConfigGetRO(ConfigTargetDir), "lastsnap.webp"), nil
		default:
			return filepath.Join(ConfigGetRO(ConfigTargetDir), "lastsnap.jpg"), nil
		}
	})

	return snapFile.(string), err
}
