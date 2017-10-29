package ircbot

import (
	"crypto/tls"
	"errors"
	"log"

	"github.com/thoj/go-ircevent"
)

type Ircbot struct {
	server   string
	tls      bool
	con      *irc.Connection
	Commands []Command
	admins   []string
}

type Command interface {
	Name() string
	Match(*irc.Event) bool
	HandleMessage(*irc.Event, *Ircbot)
	Usage() string
}

func New(nick string, srv string, enableTls bool) *Ircbot {
	ib := Ircbot{
		tls:    enableTls,
		server: srv,
	}

	ib.con = irc.IRC(nick, nick)
	ib.con.UseTLS = ib.tls
	ib.con.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	ib.con.AddCallback("PRIVMSG", ib.handleMessage)
	ib.Init()
	return &ib
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

func (ib *Ircbot) AddAdmin(host string) {
	for _, h := range ib.admins {
		if h == host {
			return
		}
	}

	ib.admins = append(ib.admins, host)
}

func (ib *Ircbot) GetAdmins() []string {
	return ib.admins
}

func (ib *Ircbot) SendMessage(c string, msg string) {
	ib.con.Privmsg(c, msg)
}

func (ib *Ircbot) Init() {
	err := ib.con.Connect(ib.server)
	if err != nil {
		log.Printf("[ircbot] Failed to connect to irc: %v\n", err)
	}
	go ib.con.Loop()
}
