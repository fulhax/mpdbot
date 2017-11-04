package irccmd

import (
	"fmt"
	"regexp"

	"github.com/fulhax/mpdbot/mpd/statistics"
	"github.com/rendom/ircFormat"
	"github.com/rendom/ircbot"
	irc "github.com/thoj/go-ircevent"
)

type UserTop struct {
	storage statistics.Storage
}

func NewUserTop(s statistics.Storage) *UserTop {
	return &UserTop{s}
}

func (i *UserTop) Name() string {
	return "UserTop"
}

func (i *UserTop) Usage() string {
	return "!top <user> - Top 10 songs for user\n!top - Top 10 songs"
}

func (i *UserTop) Match(e *irc.Event) bool {
	m, err := regexp.MatchString("^!top", e.Message())

	if err != nil {
		return false
	}

	return m
}

func (i *UserTop) HandleMessage(ev *irc.Event, ib *ircbot.Ircbot) {
	re := regexp.MustCompile("^!top (.+)")
	m := re.FindStringSubmatch(ev.Message())

	var songs []statistics.SongStats
	var err error
	if len(m) == 0 {
		songs, err = i.storage.GetTop(5)
	} else {
		songs, err = i.storage.GetUserTop(m[1], 5)
	}

	if err != nil {
		return
	}

	if len(songs) == 0 {
		ib.SendMessage(
			ev.Arguments[0],
			ircFormat.Colorize("No songs found", ircFormat.Red, ircFormat.None),
		)

		return
	}

	for _, v := range songs {
		title := ""
		if v.Song.Title != "" {
			title = v.Song.Title
		} else {
			title = v.Song.File
		}

		msg := fmt.Sprintf("%s queued %d times", title, v.TimesQueued)

		ib.SendMessage(
			ev.Arguments[0],
			msg,
		)
	}
}
