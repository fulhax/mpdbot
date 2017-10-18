package main

import (
	"regexp"

	"github.com/fulhax/mpdbot/ircbot"
	"github.com/fulhax/mpdbot/mpd"
	"github.com/rendom/ircFormat"
	irc "github.com/thoj/go-ircevent"
)

type IrcMpdUpdate struct {
	mpdClient *mpd.MpdClient
}

func (i *IrcMpdUpdate) Name() string {
	return "IrcMpdUpdate"
}

func (i *IrcMpdUpdate) Usage() string {
	return "!mpd update - updates mpd database"
}

func (i *IrcMpdUpdate) Match(e *irc.Event) bool {
	m, err := regexp.MatchString("^!mpd update", e.Message())

	if err != nil {
		return false
	}

	return m
}

func (i *IrcMpdUpdate) HandleMessage(ev *irc.Event, ib *ircbot.Ircbot) {
	err := i.mpdClient.Update()
	if err != nil {
		ib.SendMessage(
			ev.Arguments[0],
			ircFormat.Colorize("Update failed :/", ircFormat.Red, ircFormat.None),
		)
	} else {
		ib.SendMessage(
			ev.Arguments[0],
			ircFormat.Colorize("MPD Updated :D", ircFormat.Green, ircFormat.None),
		)
	}
}
