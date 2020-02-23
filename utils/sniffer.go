package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/go-httpproxy/httpproxy"
)

var (
	History    map[string]*http.Request = make(map[string]*http.Request)
	HistoryIds []string                 = []string{}
	WaitList   []string                 = []string{}
)

func SetWaitList(domain string) {
	WaitList = append(WaitList, domain)
}

func OnError(ctx *httpproxy.Context, where string,
	err *httpproxy.Error, opErr error) {
	// Log errors.
	log.Printf("ERR: %s: %s [%s]", where, err, opErr)
}

func OnAccept(ctx *httpproxy.Context, w http.ResponseWriter, r *http.Request) bool {
	// Handle local request has path "/info"
	if r.Method == "GET" && !r.URL.IsAbs() {
		if r.URL.Path == "/info" {
			w.Write([]byte("This is go-httpproxy."))
			return true
		} else if r.URL.Path == "/hist" {
			res := ""
			for no, k := range HistoryIds {
				res += fmt.Sprintf("[%d] %s\n", no, k)
			}
			w.Write([]byte(strings.TrimSpace(res)))
			return true
		} else if r.URL.Path == "/wait" {
			if s := r.FormValue("url"); s != "" {
				WaitList = append(WaitList, s)

				log.Println("Add Wait list:", s)
				w.Write([]byte("Add ok"))
				return true
			}
		} else if r.URL.Path == "/get" {
			if s := r.FormValue("url"); s != "" {
				if !strings.HasPrefix(s, "http") {
					if id, err := strconv.Atoi(s); err == nil {
						s = HistoryIds[id]
					}
				}
				if searchRq, ok := History[s]; ok {
					buf, _ := httputil.DumpRequest(searchRq, false)
					w.Write(buf)
					return true
				}

			}
		}
	}

	return false
}

func OnAuth(ctx *httpproxy.Context, authType string, user string, pass string) bool {
	// Auth test user.
	if user == "test" && pass == "1234" {
		return true
	}
	return false
}

func OnConnect(ctx *httpproxy.Context, host string) (ConnectAction httpproxy.ConnectAction, newHost string) {
	// Apply "Man in the Middle" to all ssl connections. Never change host.
	return httpproxy.ConnectMitm, host
}

func OnRequest(ctx *httpproxy.Context, req *http.Request) (
	resp *http.Response) {
	// Log proxying requests.
	log.Printf("INFO: Proxy: %s %s", req.Method, req.URL.String())
	History[req.URL.String()] = req
	HistoryIds = append(HistoryIds, req.URL.String())
	// for _, v := range WaitList {
	// 	if strings.Contains(req.Host, v) {

	// 		History[req.URL.String()] = req
	// 		HistoryIds = append(HistoryIds, req.URL.String())
	// 		break
	// 	}
	// }
	return
}

func OnResponse(ctx *httpproxy.Context, req *http.Request,
	resp *http.Response) {
	// Add header "Via: go-httpproxy".
	resp.Header.Add("Via", "go-httpproxy")
}

func RunProxyServer(port string) {
	// Create a new proxy with default certificate pair.
	prx, _ := httpproxy.NewProxy()

	// Set handlers.
	prx.OnError = OnError
	prx.OnAccept = OnAccept
	// prx.OnAuth = OnAuth
	prx.OnConnect = OnConnect
	prx.OnRequest = OnRequest
	prx.OnResponse = OnResponse

	// Listen...
	fmt.Println("[Start proxy in :8089]")
	http.ListenAndServe("localhost:"+port, prx)
	return
}

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
