package api

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	"../motion"
)

func bridgeStream(w io.Writer) bool {
	resp, err := http.Get(motion.GetStreamBaseURL())
	boundary := motion.MotionStreamBoundary

	if err != nil {
		return false
	}

	mr := multipart.NewReader(resp.Body, boundary)
	p, err := mr.NextPart()
	if err == io.EOF {
		return false
	} else if err != nil {
		return false
	}

	stream, err := ioutil.ReadAll(p)
	if err != nil {
		return false
	}

	bodyLength := strconv.Itoa(len(stream))

	stream = append([]byte("--"+boundary+"\r\nContent-type: image/jpeg\r\nContent-Length:"+bodyLength+"\r\n\r\n"), stream...)
	stream = append(stream, []byte("\r\n")...)

	_, err = w.Write(stream)
	if err != nil {
		return false
	}

	return true
}
