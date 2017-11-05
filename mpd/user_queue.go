package mpd

import (
	"fmt"
	"math/rand"
	"time"
)

type userQueue struct {
	queue   []queueItem
	autodj  bool
	user    string
	handler *QueueHandler
}

func (u *userQueue) addToQueue(item queueItem) {
	u.queue = append(u.queue, item)
}

func (u *userQueue) getRandomSong() (queueItem, error) {
	songs, err := u.handler.statsStorage.GetUserTop(u.user, 200)
	if err != nil {
		return queueItem{}, err
	}

	if len(songs) == 0 {
		return queueItem{}, fmt.Errorf("No songs found")
	}

	slen := len(songs)
	k := 0
	if slen > 1 {
		rand.Seed(time.Now().Unix())
		k = rand.Intn(len(songs) - 1)
	}

	return queueItem{
		u.user,
		songs[k].Song.File,
		songs[k].Song.Title,
		time.Now(),
	}, nil
}

func (u *userQueue) EnableAutoDj() error {
	songs, _ := u.handler.statsStorage.GetUserTop(u.user, 15)
	if len(songs) < 15 {
		return fmt.Errorf("Requires 15 songs for autodj")
	}

	u.autodj = true
	return nil
}

func (u *userQueue) DisableAutoDj() {
	u.autodj = false
}

func (u *userQueue) pullNextSong() queueItem {
	if len(u.queue) == 0 {
		if u.autodj {
			s, _ := u.getRandomSong()
			return s
		} else {
			return queueItem{}
		}
	}

	qi := u.queue[0]
	u.queue = append(u.queue[:0], u.queue[1:]...)
	return qi
}
