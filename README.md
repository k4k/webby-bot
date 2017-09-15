# Webby Bot
Webby is an IRC bot written in Go. The initial goal is to provide helpers for
parsing URLs to provide the URL title and a shortend URL.

# Building

Building requires go dep:
`$ go get -u github.com/golang/dep/cmd/dep`


```
$ git clone git@github.com:k4k/webby-bot.git
$ cd webby-bot
$ dep ensure
$ go build
```

# Using Webby

Change these values before you build:

const server = "irc.freenode.net" // IRC Server                                 
const port = 7000                 // IRC Server Port                            
const channel = "#myircchannel"   // IRC Channel                                
const nick = "webby-urlbot"       // IRC Nick                                   
const username = "webby-urlbot"   // IRC Username                               
const realname = "Webby Urlbot"   // IRC Real Name```
```
