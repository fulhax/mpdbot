package main

import (
	"fmt"
	"regexp"

	"github.com/fulhax/mpdbot/ircbot"
	"github.com/fulhax/mpdbot/mpd"
	"github.com/rendom/ircFormat"
	irc "github.com/thoj/go-ircevent"
)

type IrcAddSong struct {
	mpdClient *mpd.MpdClient
}

func (i *IrcAddSong) Name() string {
	return "IrcAddSong"
}

func (i *IrcAddSong) Usage() string {
	return "!mpd add <searchstr> - Search for song in library and adds it to queue."
}

func (i *IrcAddSong) Match(e *irc.Event) bool {
	m, err := regexp.MatchString("^!mpd add .+", e.Message())

	if err != nil {
		return false
	}

	return m
}

func (i *IrcAddSong) HandleMessage(ev *irc.Event, ib *ircbot.Ircbot) {
	re := regexp.MustCompile("!mpd add (.+)")
	m := re.FindStringSubmatch(ev.Message())

	if len(m) == 0 {
		return
	}

	file, err := queueHandler.AddToQueue(ev.Nick, m[1])
	if err != nil {
		return
	}

	if file.File != "" {
		msg := fmt.Sprintf("Added %s to queue", file.File)
		ib.SendMessage(
			ev.Arguments[0],
			ircFormat.Colorize(msg, ircFormat.Green, ircFormat.None),
		)
	} else {
		msg := fmt.Sprintf("Unable to find song %s", m[1])
		ib.SendMessage(
			ev.Arguments[0],
			ircFormat.Colorize(msg, ircFormat.Red, ircFormat.None),
		)
	}
}
