package utils

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type ProxyCli string

func (cli *ProxyCli) Hist() (out string) {
	res, err := http.Get("http://localhost:8089/hist")
	if err != nil {
		log.Fatal(err)
	}
	buf, _ := ioutil.ReadAll(res.Body)
	out = string(buf)
	return
}

func (cli *ProxyCli) Get(url string) (out string) {
	res, err := http.Get("http://localhost:8089/get?url=" + url)
	if err != nil {
		log.Fatal(err)
	}
	buf, _ := ioutil.ReadAll(res.Body)
	out = string(buf)
	return
}

func (cli *ProxyCli) AddCheck(url string) (out string) {
	res, err := http.Get("http://localhost:8089/wait?url=" + url)
	if err != nil {
		log.Fatal(err)
	}
	buf, _ := ioutil.ReadAll(res.Body)
	out = string(buf)
	return
}

func SimpleGet(url string) (output string) {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   12 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	outbuf, _ := httputil.DumpResponse(res, false)
	output += strings.TrimSpace(string(outbuf))
	buf, _ := ioutil.ReadAll(res.Body)
	output += "\r\n\r\n" + strings.TrimSpace(string(buf))
	return

}
