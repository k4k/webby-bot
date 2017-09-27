// Copyright [2017] [Ted Wood <ted@xy0.org>]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/thoj/go-ircevent"
	"github.com/zpnk/go-bitly"
	"golang.org/x/net/html"
)

// Define your settings here.
// TODO: Replace with config and cmd args using cobra
const server = "irc.freenode.net" // IRC Server
const port = 7000                 // IRC Server Port
const channel = "#myircchannel"   // IRC Channel
const nick = "webby-urlbot"       // IRC Nick
const username = "webby-urlbot"   // IRC Username
const realname = "Webby Urlbot"   // IRC Real Name
const token = "<token>"           // bitly OAuth Token

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
			shortURL := shortenURL(urlstring)
			irccon.Privmsg(channel, fmt.Sprintf("%s -- %s", shortURL, title))
		}
	})

	//irccon.AddCallback("336", func(e *irc.Event) {})
	err := irccon.Connect(fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		fmt.Errorf("Err %s", err)
	}
	irccon.Loop()
}

// Bitly Functions
func shortenURL(longurl string) string {
	b := bitly.New(token)
	shortURL, shortErr := b.Links.Shorten(longurl)
	if shortErr != nil {
		fmt.Errorf("Error: %s", shortErr)
	}
	return shortURL.URL
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
