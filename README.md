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
