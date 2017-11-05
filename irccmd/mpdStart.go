package irccmd

import (
	"regexp"

	"github.com/fulhax/mpdbot/mpd"
	"github.com/rendom/ircbot"
	irc "github.com/thoj/go-ircevent"
)

type MpdStart struct {
	queueHandler *mpd.QueueHandler
}

func NewMpdStart(q *mpd.QueueHandler) *MpdStart {
	return &MpdStart{q}
}

func (i *MpdStart) Name() string {
	return "MpdStart"
}

func (i *MpdStart) Usage() string {
	return "!mpd start - Start mpdbot queue"
}

func (i *MpdStart) Match(e *irc.Event) bool {
	m, err := regexp.MatchString("^!mpd start", e.Message())

	if err != nil {
		return false
	}

	return m
}

func (i *MpdStart) HandleMessage(ev *irc.Event, ib *ircbot.Ircbot) {
	i.queueHandler.StartQueue()
}
