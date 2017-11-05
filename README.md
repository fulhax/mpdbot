# mpdbot 
[![Build Status](https://travis-ci.org/fulhax/mpdbot.svg?branch=master)](https://travis-ci.org/fulhax/mpdbot)

Install
`go install github.com/fulhax/mpdbot/cmd/mpdbot`

```
Usage of ./mpdbot:
      --debug                Enable debug mode
      --httpPort string      Http port (default "8888")
      --ircEnabled           Enable irc bot (default true)
      --ircNick string       Irc nick (default "mpdbot")
      --ircServer string     irc server (default "127.0.0.1:6697")
      --ircTls               irc tls (default true)
      --mpd string           mpd host (default "127.0.0.1:6600")
      --mpdPassword string   mpd password
```
#### Irc commands
| Event | |
| --- | --- |
| !help | List all available commands |
| !np | Now playing |
| !mpd update | Updates mpd index |
| !mpd add <search> | Add song to queue (fuzzy search) |
| !top | Top 5 queued songs |
| !top <user> | Top 5 queued song by user |
| !autodj | Enable autodj (If user queue is empty it will fetch random song from he's top 200) |
      
#### HTTP api
| URI | Method | Params |  |
| --- | --- | --- | --- |
| /current | GET | | |
| /next | POST | | |
| /add | GET | user, song |  |
| /search| GET | search |  |
| /status | GET | |  |
| /top | GET | user(optional) |  |


config.yml 
```
debug: false
mpd: 127.0.0.1:6600
mpdPassword: password
httpPort: 8080
ircServer: 127.0.0.1:6697
ircTls: true
ircEnabled: true
ircNick: "mpdbot"
```

