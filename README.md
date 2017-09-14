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

Right now, the server and channel to join need to be statically defined before
you build. I have plans to support a config file and command line flags for these
down the road (maybe).
