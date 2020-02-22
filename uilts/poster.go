package utils

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

func SendHTTPFromReader(r io.Reader) (output string, err error) {
	buf := bufio.NewReader(r)
	req := new(http.Request)
	req, err = http.ReadRequest(buf)
	if err == io.EOF {
		return
	}
	client := http.Client{
		Timeout: 7 * time.Second,
	}
	if res, ierr := client.Do(req); err != nil {
		return "", ierr
	} else {
		var buf []byte
		buf, err = httputil.DumpResponse(res, true)
		if err != nil {
			return
		}
		output = string(buf)
	}
	return
}
