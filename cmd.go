package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/Qingluan/VimExtend/utils"
)

var (
	usereq      bool
	proxyserver bool
	query       string
	getreq      string
	logreq      string
	urlget      string
	hist        bool
)

func main() {
	flag.BoolVar(&usereq, "r", true, "set true to read stdin to parse req then make request")
	flag.StringVar(&query, "q", "a", "true to cssselect content from stdin")
	flag.StringVar(&getreq, "g", "", "get url from proxy")
	flag.StringVar(&logreq, "l", "", "log req by this domain")
	flag.StringVar(&urlget, "u", "", "req by this url")
	flag.BoolVar(&hist, "ls", false, "show hist in proxy server")
	flag.BoolVar(&proxyserver, "S", false, "set Server start")

	flag.Parse()
	if proxyserver {
		utils.RunProxyServer("8089")
		os.Exit(0)
	}
	if urlget != "" {
		fmt.Println(utils.SimpleGet(urlget))
		os.Exit(0)
	}
	if hist {
		c := new(utils.ProxyCli)
		fmt.Println(c.Hist())
		os.Exit(0)
	}
	if logreq != "" {
		c := new(utils.ProxyCli)
		fmt.Println(c.AddCheck(logreq))
		os.Exit(0)
	}
	if getreq != "" {
		c := new(utils.ProxyCli)
		fmt.Println(c.Get(getreq))
		os.Exit(0)
	}

	if usereq {
		reader := bufio.NewReader(os.Stdin)
		// var output []rune
		buffer := bytes.NewBuffer([]byte{})
		buf := make([]byte, 4096)
		io.CopyBuffer(buffer, reader, buf)

		if header, body, err := utils.SendHTTPFromReader(buffer); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(header)
			fmt.Println("\r\n\r\n")
			fmt.Println(body)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		// var output []rune
		tmp, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println(err.Error())
		}
		pts := bytes.SplitN(tmp, []byte("<html"), 2)
		html := "<html" + string(pts[1])
		buffer := bytes.NewBuffer([]byte(html))
		doc, err := goquery.NewDocumentFromReader(buffer)
		if err != nil {
			fmt.Println(err.Error())
		}
		doc.Find(query).Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			if query == "a" {
				if href, ok := s.Attr("href"); ok {
					fmt.Println(href)
				}
			} else {
				if content, err := s.Html(); err == nil {
					fmt.Println(content)
				}
			}
		})
	}

	// for _, v := range strings.Split(buffer.String(), "\n") {
	// 	fmt.Println("-> ", v)
	// }
}
