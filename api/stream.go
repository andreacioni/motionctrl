package api

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"../motion"
)

func bridgeStream(w io.Writer) bool {
	resp, err := http.Get(motion.GetStreamBaseURL())

	if err != nil {
		return false
	}

	mr := multipart.NewReader(resp.Body, motion.MotionStreamBoundary)
	p, err := mr.NextPart()
	if err == io.EOF {
		return false
	}
	if err != nil {
		return false
	}
	stream, err := ioutil.ReadAll(p)
	if err != nil {
		return false
	}

	_, err = w.Write(stream)

	if err != nil {
		return false
	}

	return true
}
