package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/thoj/go-ircevent"
	"golang.org/x/net/html"
)

// Define your settings here.
// TODO: Replace with config and cmd args using cobra
const server = "irc.freenode.net" // IRC Server
const port = 7000                 // IRC Server Port
const channel = "#myircchannel"   // IRC Channel
const nick = "urlbot"             // IRC Nick
const username = "fooas-urlbot"   // IRC Username
const realname = "FOAAS Urlbot"   // IRC Real Name

// IRC Functions
func ircConnect(srv string, port int, ch string) {
	irccon := irc.IRC(nick, "FOAAS")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {

		URL, urlstring := isURL(e.Message())
		if URL {
			title := getTitle(urlstring)
			irccon.Privmsg(channel, title)
		}
	})

	//irccon.AddCallback("336", func(e *irc.Event) {})
	err := irccon.Connect(fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		fmt.Printf("Err %s", err)
	}
	irccon.Loop()
}

// HTML Functions
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

func getTitle(gettitle string) string {

	resp, getErr := http.Get(gettitle)
	var title string
	if getErr != nil {
		panic(getErr)
	}
	defer resp.Body.Close()

	title, _ = GetHtmlTitle(resp.Body)
	if title == "" {
		title = "No title found"
	}
	return title
}

func isURL(checkurl string) (bool, string) {
	var realURL bool = false
	var realItem string = ""
	urlslice := strings.Split(checkurl, " ")
	for _, item := range urlslice {
		u, err := url.Parse(item)
		if err == nil {
			if u.Scheme == "http" || u.Scheme == "https" {
				realURL = true
				realItem = item
				break
			}
		}
	}
	return realURL, realItem
}

func main() {
	ircConnect(server, port, channel)
}
