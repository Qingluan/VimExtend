package utils

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	// "net/http/httputils"
)

func SendHTTPFromReader(r io.Reader) (header string, body string, err error) {
	buf := bufio.NewReader(r)
	req := new(http.Request)
	var reqR *http.Request
	reqR, err = http.ReadRequest(buf)

	// u.Path = path.Join(u.Path,)
	req, _ = http.NewRequest(reqR.Method, reqR.URL.String(), nil)
	req.PostForm = reqR.PostForm
	req.Header = reqR.Header
	client := new(http.Client)
	var outbuf []byte
	if err == io.EOF {
		return
	} else if err != nil {
		return
	}
	res, ierr := client.Do(req)
	if ierr != nil {
		return "", "", ierr
	}
	defer res.Body.Close()
	outbuf, _ = httputil.DumpResponse(res, false)
	outbuf2, _ := ioutil.ReadAll(res.Body)
	// buf, ierr := httputil.DumpResponse(res, true)
	if err != nil {
		return
	}
	header = string(outbuf)
	body = string(outbuf2)
	return
}
