package ircbot

import (
	"errors"
	"fmt"

	"github.com/thoj/go-ircevent"
)

type Ircbot struct {
	server   string
	tls      bool
	con      *irc.Connection
	Commands []Command
}

type Command interface {
	Name() string
	Match(*irc.Event) bool
	HandleMessage(*irc.Event, *Ircbot)
	Usage() string
}

func NewIrcBot(nick string, srv string, tls bool) *Ircbot {
	var ib *Ircbot
	ib.tls = tls
	ib.server = srv

	ib.con = irc.IRC(nick, nick)
	ib.con.UseTLS = true
	ib.con.AddCallback("PRIVMSG", ib.handleMessage)
	ib.Init()
	return ib
}

func (ib *Ircbot) handleMessage(e *irc.Event) {
	if e.Nick == ib.con.GetNick() {
		return
	}

	for _, cmd := range ib.Commands {
		if cmd.Match(e) {
			cmd.HandleMessage(e, ib)
		}
	}
}

func (ib *Ircbot) AddCommand(cmd Command) error {
	for _, c := range ib.Commands {
		if c.Name() == cmd.Name() {
			return errors.New("Command already registered")
		}
	}

	ib.Commands = append(ib.Commands, cmd)
	return nil
}

func (ib *Ircbot) RemoveCommand(name string) error {
	for i, c := range ib.Commands {
		if c.Name() == name {
			ib.Commands = append(ib.Commands[:i], ib.Commands[i+1:]...)
			return nil
		}
	}

	return errors.New("Command not found")
}

func (ib *Ircbot) SendMessage(c string, msg string) {
	ib.con.Privmsg(c, msg)
}

func (ib *Ircbot) Init() {
	ib.AddCommand(Usage{})

	err := ib.con.Connect(ib.server)
	if err != nil {
		fmt.Printf("[ircbot] Failed to connect to irc:  %s", err)
	}
	ib.con.Loop()
}
