package mpd

import (
	"fmt"
	"log"
	"time"

	"github.com/fulhax/mpdbot/mpd/statistics"

	"github.com/fhs/gompd/mpd"
)

type QueueHandler struct {
	Client       client
	StatsStorage statistics.Storage
	currentUser  string
	usersQueues  []*userQueue
}

type userQueue struct {
	queue  []queueItem
	autodj bool
	user   string
}

func (u *userQueue) addToQueue(item queueItem) {
	u.queue = append(u.queue, item)
}
func (u *userQueue) pullNextSong() queueItem {
	if len(u.queue) == 0 {
		// TODO if autodj get random song from statistics
		return queueItem{}
	}

	qi := u.queue[0]
	u.queue = append(u.queue[:0], u.queue[1:]...)
	return qi
}

type queueItem struct {
	User  string
	File  string
	Title string
	Added time.Time
}

func (q *QueueHandler) Init(addr, pw string) (err error) {
	go q.handlePlaylist(addr, pw)
	return nil
}

func (q *QueueHandler) songInQueue(file string) bool {
	for _, uq := range q.usersQueues {
		for _, v := range uq.queue {
			if v.File == file {
				return true
			}

		}
	}

	return false
}

func (q *QueueHandler) getUserQueue(u string) *userQueue {
	for _, uq := range q.usersQueues {
		if uq.user == u {
			return uq
		}
	}

	uq := &userQueue{autodj: false, user: u}
	q.usersQueues = append(q.usersQueues, uq)

	return uq
}

// AddToQueue adds song to user queue, will return error if already queued
func (q *QueueHandler) AddToQueue(user string, title string, file string) (queueItem, error) {
	if q.songInQueue(file) {
		return queueItem{}, fmt.Errorf("%s is already queued", file)
	}

	item := queueItem{
		User:  user,
		File:  file,
		Title: title,
		Added: time.Now(),
	}

	q.getUserQueue(user).addToQueue(item)
	err := q.StatsStorage.AddSong(item.Title, item.File, user)
	if err != nil {
		log.Printf("Error while saving statistics: %v\n", err)
	}

	return item, nil
}

func (q *QueueHandler) pullNextSong() (file queueItem, err error) {
	picked := -1
	var qi queueItem
	for i, uq := range q.usersQueues {
		qi = uq.pullNextSong()
		if qi.File != "" {
			picked = i
			break
		}
	}

	len := len(q.usersQueues)
	if picked != -1 && len > 1 {
		m := q.usersQueues[picked]
		q.currentUser = m.user
		q.usersQueues = append(q.usersQueues[:picked], q.usersQueues[picked+1:]...)
		q.usersQueues = append(q.usersQueues, m)
	}

	return qi, nil
}

// Handle mpd playlist and add new songs from queue
func (q *QueueHandler) handlePlaylist(addr, pw string) (err error) {
	w, err := mpd.NewWatcher("tcp", addr, pw)
	if err != nil {
		return err
	}
	defer w.Close()

	go func() {
		for err := range w.Error {
			fmt.Println("MPD watcher error: ", err)
		}
	}()

	for range w.Event {
		status, _ := q.Client.GetStatus()
		if status.State == "stop" {
			q.Client.ClearPlaylist()
			q.queueNextSong()
			q.Client.Play(0)
		}
	}

	return nil
}

func (q *QueueHandler) queueNextSong() string {
	next, _ := q.pullNextSong()

	if next.File != "" {
		q.Client.AddSong(next.File)
		return next.File
	} else {
		song, err := q.Client.GetRandomSong()
		if err == nil {
			q.Client.AddSong(song)
		}

		return song
	}
}
