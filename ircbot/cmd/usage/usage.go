package usage

import (
	"fmt"
	"regexp"

	irc "github.com/thoj/go-ircevent"
)

type Usage struct{}

func (usage *Usage) Name() string {
	return "Usage"
}

func (usage *Usage) Usage() string {
	return "!help list all registered commands"
}

func (usage *Usage) Match(*irc.Event) bool {
	m, err := regexp.MatchString("^!help", msg)

	if err != nil {
		return false
	}

	return m
}

func (usage *Usage) HandleMessage(ev *irc.Event, ib *Ircbot) {
	// list all commands usage
	for _, v := range ib.Commands {
		fmt.Println("lol")
	}
}
