package motion

import (
	"fmt"
	"github.com/andreacioni/motionctrl/utils"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func Snapshot() error {
	_, err := webControlGet("/action/snapshot", func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(SnapshotDetectionRegex, body) {
			return "", fmt.Errorf("unable to take snapshot (%s)", body)
		}

		return "", nil
	})

	return err
}

func Capture() (string, error) {
	resp, err := http.Get(GetStreamBaseURL())
	boundary := "BoundaryString"

	if err != nil {
		return "", err
	}

	mr := multipart.NewReader(resp.Body, boundary)
	p, err := mr.NextPart()
	if err != nil {
		return "", err
	}

	stream, err := ioutil.ReadAll(p)
	if err != nil {
		return "", err
	}

	tempFile := filepath.Join(os.TempDir(), "capture.jpg")

	if err := ioutil.WriteFile(tempFile, stream, 0600); err != nil {
		return "", err
	}

	return tempFile, nil
}
