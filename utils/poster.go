package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func SendHTTPFromReader(r io.Reader) (output string, err error) {
	buf := bufio.NewReader(r)
	req := new(http.Request)
	req, err = http.ReadRequest(buf)
	var outbuf []byte
	if err == io.EOF {
		return
	} else if err != nil {
		return
	}
	req.RequestURI = ""
	fmt.Println("read header ok")
	client := http.Client{
		Timeout: 7 * time.Second,
	}
	res, ierr := client.Do(req)
	if ierr != nil {
		return "", ierr
	}
	defer res.Body.Close()
	outbuf, err = ioutil.ReadAll(res.Body)
	// buf, ierr := httputil.DumpResponse(res, true)
	if err != nil {
		return
	}
	output = string(outbuf)

	return
}
