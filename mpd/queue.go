package mpd

import (
	"fmt"
	"time"

	"github.com/fhs/gompd/mpd"
)

type QueueHandler struct {
	MpdClient   *MpdClient
	currentSong string
	currentUser string
	lastQueued  string
	queue       []QueueItem
}

type QueueItem struct {
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
	for _, v := range q.queue {
		if v.File == file {
			return true
		}
	}

	return false
}

func (q *QueueHandler) AddToQueue(user string, song string) (QueueItem, error) {
	sr, err := q.MpdClient.SearchInLibrary(song)
	if err != nil {
		return QueueItem{}, err
	}

	if len(sr) > 0 {
		if q.songInQueue(sr[0].File) == false {
			item := QueueItem{
				User:  user,
				File:  sr[0].File,
				Title: sr[0].Title,
				Added: time.Now(),
			}
			q.queue = append(q.queue, item)
			return item, nil
		}
	}

	return QueueItem{}, nil
}

func (q *QueueHandler) pullNextSong() (file QueueItem, err error) {
	idx := 0
	for k, i := range q.queue {
		if file.File == "" || i.User != file.User {
			file = i
			idx = k
		}

		if q.currentUser != i.User {
			break
		}
	}
	if file.File != "" {
		q.queue = append(q.queue[:idx], q.queue[idx+1:]...)
	}
	return file, nil
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
		q.currentUser = next.User
	} else {
		song, err := q.MpdClient.GetRandomSong()
		if err == nil {
			q.MpdClient.AddSong(song)
		}
	}

	return
}
