package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-httpproxy/httpproxy"
)

var defaultcacert = []byte(`-----BEGIN CERTIFICATE-----
MIICsjCCAjigAwIBAgIUQUrVa2KWAoSCIx5PcIgwZR/g41QwCgYIKoZIzj0EAwIw
gY8xCzAJBgNVBAYTAkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYDVQQHDAdCZWlq
aW5nMRUwEwYDVQQKDAxTY2huaXNzVmllbGUxGDAWBgNVBAsMD1pob25nIEd1YW4g
Q2h1bjELMAkGA1UEAwwCQ0sxHjAcBgkqhkiG9w0BCQEWD2NoZWNrQGdtYWlsLmNv
bTAeFw0yMDAyMjMxMDQ0NTZaFw0zMDAyMjAxMDQ0NTZaMIGPMQswCQYDVQQGEwJD
TjEQMA4GA1UECAwHQmVpamluZzEQMA4GA1UEBwwHQmVpamluZzEVMBMGA1UECgwM
U2Nobmlzc1ZpZWxlMRgwFgYDVQQLDA9aaG9uZyBHdWFuIENodW4xCzAJBgNVBAMM
AkNLMR4wHAYJKoZIhvcNAQkBFg9jaGVja0BnbWFpbC5jb20wdjAQBgcqhkjOPQIB
BgUrgQQAIgNiAASThMqDNlhhd12twi9O9zs67Yqci8qF2sSL5HToBiQcmk/gp6GL
zNch3UkBvfNoJzlpqTrzTd13uZXbhAG2FsD98u9EvIjhoFjWvJ63qAgnkIVMMJKC
5rPBME+D2DCCUiijUzBRMB0GA1UdDgQWBBRZT2B3Crc707bantz4pEUwi8jmzDAf
BgNVHSMEGDAWgBRZT2B3Crc707bantz4pEUwi8jmzDAPBgNVHRMBAf8EBTADAQH/
MAoGCCqGSM49BAMCA2gAMGUCMQCMskLiTX9rgGLJnhjRUlfjZ9n691wViZ4sZXsq
/jn63C5hvTuIchw4PMb6SCvWAysCMBB0OjpFPnElTw8daUVN4sENLEGLQtjRa72E
fyMQobglpkUsRqb12aX02Kdmxh+LtA==
-----END CERTIFICATE-----`)
var defaultcakey = []byte(`-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDBsk17Moudo/kXZZVDgVLRgNXErbsrnbvOJGDFKcQa6bqnYD0RLbEOh
+MWR3cWM99ygBwYFK4EEACKhZANiAASThMqDNlhhd12twi9O9zs67Yqci8qF2sSL
5HToBiQcmk/gp6GLzNch3UkBvfNoJzlpqTrzTd13uZXbhAG2FsD98u9EvIjhoFjW
vJ63qAgnkIVMMJKC5rPBME+D2DCCUig=
-----END EC PRIVATE KEY-----`)

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
	return true
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
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set handlers.
	prx.OnError = OnError
	prx.OnAccept = OnAccept
	// prx.OnAuth = OnAuth
	prx.OnConnect = OnConnect
	prx.OnRequest = OnRequest
	prx.OnResponse = OnResponse

	// Listen...
	fmt.Println("[Start proxy in :8089]")
	server := &http.Server{
		Addr:         "localhost:" + port,
		Handler:      prx,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	listenErrChan := make(chan error)
	go func() {
		listenErrChan <- server.ListenAndServe()
	}()
	log.Printf("Listening HTTP %s", server.Addr)

	cert, _ := tls.X509KeyPair(defaultcacert, defaultcakey)
	serverHTTPS := &http.Server{
		Addr:         "localhost:8090",
		Handler:      prx,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionSSL30,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			Certificates: []tls.Certificate{cert},
		},
	}
	listenHTTPSErrChan := make(chan error)
	go func() {
		listenHTTPSErrChan <- serverHTTPS.ListenAndServeTLS("", "")
	}()
	log.Printf("Listening HTTPS %s", serverHTTPS.Addr)
mainloop:
	for {
		select {
		case <-sigChan:
			break mainloop
		case listenErr := <-listenErrChan:
			if listenErr != nil && listenErr == http.ErrServerClosed {
				break mainloop
			}
			log.Fatal(listenErr)
		case listenErr := <-listenHTTPSErrChan:
			if listenErr != nil && listenErr == http.ErrServerClosed {
				break mainloop
			}
			log.Fatal(listenErr)
		}
	}

	shutdown := func(srv *http.Server, wg *sync.WaitGroup) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err == context.DeadlineExceeded {
			log.Printf("Force shutdown %s", srv.Addr)
		} else {
			log.Printf("Graceful shutdown %s", srv.Addr)
		}
		wg.Done()
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go shutdown(server, wg)
	wg.Add(1)
	go shutdown(serverHTTPS, wg)
	wg.Wait()

	log.Println("Finished")
	// http.ListenAndServe("localhost:"+port, prx, make(map[string]func(*http.Server, *tls.Conn, http.Handler)))
	return
}
