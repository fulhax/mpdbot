package irccmd

import (
	"fmt"
	"regexp"

	"github.com/fulhax/mpdbot/mpd"
	"github.com/rendom/ircFormat"
	"github.com/rendom/ircbot"
	irc "github.com/thoj/go-ircevent"
)

type Autodj struct {
	mpdClient    *mpd.MpdClient
	queueHandler *mpd.QueueHandler
}

func NewAutodj(m *mpd.MpdClient, q *mpd.QueueHandler) *Autodj {
	return &Autodj{m, q}
}

func (i *Autodj) Name() string {
	return "Autodj"
}

func (i *Autodj) Usage() string {
	return "!autodj - Toggle autodj."
}

func (i *Autodj) Match(e *irc.Event) bool {
	m, err := regexp.MatchString("^!autodj", e.Message())

	if err != nil {
		return false
	}

	return m
}

func (i *Autodj) HandleMessage(ev *irc.Event, ib *ircbot.Ircbot) {
	v, err := i.queueHandler.ToggleAutodj(ev.Nick)

	if err != nil {
		ib.SendMessage(
			ev.Arguments[0],
			ircFormat.Colorize(err.Error(), ircFormat.Red, ircFormat.None),
		)
		return
	}

	msg := fmt.Sprintf("Autodj: %v", v)
	ib.SendMessage(
		ev.Arguments[0],
		ircFormat.Colorize(msg, ircFormat.Green, ircFormat.None),
	)
}
