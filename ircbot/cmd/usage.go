package cmd

import (
	"fmt"
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
	fmt.Println("hitting match")
	m, err := regexp.MatchString("^!lol", e.Message())

	if err != nil {
		fmt.Println("no match")
		return false
	}

	fmt.Println("match stuff")
	return m
}

func (usage *Usage) HandleMessage(ev *irc.Event, ib *ircbot.Ircbot) {
	fmt.Println("Handle message")
	// list all commands usage
	for _, v := range ib.Commands {
		ib.SendMessage(ev.Arguments[0], v.Usage())
		fmt.Println("lol")
	}
}
