package cmd

import (
	"regexp"

	"github.com/fulhax/mpdbot/ircbot"
	irc "github.com/thoj/go-ircevent"
)

type Usage struct{}

func (usage *Usage) Name() string {
	return "Usage"
}

func (usage *Usage) Usage() string {
	return "!lol list all registered commands"
}

func (usage *Usage) Match(e *irc.Event) bool {
	m, err := regexp.MatchString("^!lol", e.Message())

	if err != nil {
		return false
	}

	return m
}

func (usage *Usage) HandleMessage(ev *irc.Event, ib *ircbot.Ircbot) {
	// list all commands usage
	for _, v := range ib.Commands {
		ib.SendMessage(ev.Arguments[0], v.Usage())
	}
}
