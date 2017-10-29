package irccmd

import (
	"regexp"

	"github.com/fulhax/mpdbot/mpd"
	"github.com/rendom/ircbot"
	irc "github.com/thoj/go-ircevent"
)

type IrcMpdNp struct {
	mpdClient *mpd.MpdClient
}

func NewNp(m *mpd.MpdClient) *IrcMpdNp {
	return &IrcMpdNp{m}
}

func (i *IrcMpdNp) Name() string {
	return "IrcMpdNp"
}

func (i *IrcMpdNp) Usage() string {
	return "!np current song"
}

func (i *IrcMpdNp) Match(e *irc.Event) bool {
	m, err := regexp.MatchString("^!np", e.Message())

	if err != nil {
		return false
	}

	return m
}

// TODO include state in response
func (i *IrcMpdNp) HandleMessage(ev *irc.Event, ib *ircbot.Ircbot) {
	song, err := i.mpdClient.CurrentSong()
	if err == nil {
		ib.SendMessage(ev.Arguments[0], song)
	}
}
