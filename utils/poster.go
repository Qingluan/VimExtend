package utils

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	// "net/http/httputils"
	"net/url"
	"path"
)

func SendHTTPFromReader(r io.Reader) (header string, body string, err error) {
	buf := bufio.NewReader(r)
	req := new(http.Request)
	var reqR *http.Request
	reqR, err = http.ReadRequest(buf)
	u, err := url.Parse("http://" + reqR.Host)
	u.Path = path.Join(u.Path, reqR.URL.String())
	req, _ = http.NewRequest(reqR.Method, u.String(), nil)
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
