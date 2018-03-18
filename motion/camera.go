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

		if targetDir, err := ConfigGet(ConfigTargetDir); err == nil {
			switch snapExt {
			case "ppm":
				return filepath.Join(targetDir.(string), "lastsnap.ppm"), nil
			case "webp":
				return filepath.Join(targetDir.(string), "lastsnap.webp"), nil
			default:
				return filepath.Join(targetDir.(string), "lastsnap.jpg"), nil
			}
		} else {
			return nil, err
		}

	})

	return snapFile.(string), err
}
