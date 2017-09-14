package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/clarkf/libirc"
	"golang.org/x/net/html"
)

// Define you server, port and channel here
const server = "irc.freenode.net"  // IRC Server
const port = 6667                  // IRC Server Port
const channel = "#somechannelhere" // IRC Channel

func ircConnect(svr string, port int, cnl string) {
	c := libirc.NewClient("webby", "webby-bot", "Webby Bot")

	connErr := c.ConnectAndListen(fmt.Sprintf("%s:%d", svr, port))
	if connErr != nil {
		panic(connErr)
	}

	c.Join(cnl)
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}

func GetHtmlTitle(r io.Reader) (string, bool) {
	doc, docErr := html.Parse(r)
	if docErr != nil {
		panic("Failed to parse url")
	}

	return traverse(doc)
}

func main() {

	ircConnect(server, port, channel)
	time.Sleep(time.Duration(5) * time.Second)

	url := "http://google.com"
	resp, getErr := http.Get(url)
	if getErr != nil {
		panic(getErr)
	}
	defer resp.Body.Close()

	if title, ok := GetHtmlTitle(resp.Body); ok {
		println(title)
	} else {
		println("No title found")
	}
}
