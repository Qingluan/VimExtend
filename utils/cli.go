package utils

import (
	"io/ioutil"
	"log"
	"net/http"
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
