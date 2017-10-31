package mpd

import (
	"fmt"
	"github.com/fulhax/mpdbot/mpd/statistics"
	"time"

	"github.com/fhs/gompd/mpd"
	"log"
)

type QueueHandler struct {
	MpdClient    *MpdClient
	StatsStorage statistics.Storage
	currentUser  string
	usersQueues  []*userQueue
}

type userQueue struct {
	queue  []queueItem
	autodj bool
	user   string
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

func (q *QueueHandler) Init() (err error) {
	go q.handlePlaylist()
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

func (q *QueueHandler) AddToQueue(user string, song string) (queueItem, error) {
	sr, err := q.MpdClient.SearchInLibrary(song)
	if err != nil {
		return queueItem{}, err
	}

	if len(sr) > 0 {
		if q.songInQueue(sr[0].File) == false {
			item := queueItem{
				User:  user,
				File:  sr[0].File,
				Title: sr[0].Title,
				Added: time.Now(),
			}
			uq := q.getUserQueue(user)
			uq.queue = append(uq.queue, item)
			err := q.StatsStorage.AddSong(item.Title, item.File, user)
			if err != nil {
				log.Printf("Error while saving statistics: %v\n", err)
			}
			return item, nil
		}
	}

	return queueItem{}, nil
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
func (q *QueueHandler) handlePlaylist() (err error) {
	w, err := mpd.NewWatcher("tcp", q.MpdClient.addr, q.MpdClient.password)
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
		status, _ := q.MpdClient.GetStatus()
		if status.State == "stop" {
			q.MpdClient.ClearPlaylist()
			q.queueNextSong()
			q.MpdClient.Play(0)
		}
	}

	return nil
}

func (q *QueueHandler) queueNextSong() {
	next, _ := q.pullNextSong()

	if next.File != "" {
		q.MpdClient.AddSong(next.File)
	} else {
		song, err := q.MpdClient.GetRandomSong()
		if err == nil {
			q.MpdClient.AddSong(song)
		}
	}
}
