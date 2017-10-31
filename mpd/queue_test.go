package mpd

import (
	"fmt"
	"testing"
	"time"
)

func generateFakeQueue() QueueHandler {
	uq := make([]*userQueue, 3)
	for k, _ := range uq {
		u := fmt.Sprintf("user%d", k)
		uq[k] = &userQueue{
			user: u,
			queue: []queueItem{
				queueItem{File: fmt.Sprintf("%dlol.mp3", k), Added: time.Now(), User: u},
				queueItem{File: fmt.Sprintf("%dlol1.mp3", k), Added: time.Now(), User: u},
			},
		}
	}

	q := QueueHandler{
		currentUser: "user1",
		usersQueues: uq,
	}

	return q
}

func TestSongInQueue(t *testing.T) {
	q := generateFakeQueue()

	files := []struct {
		f string
		e bool
	}{
		{"0lol1.mp3", true},
		{"2lol.mp3", true},
		{"0foo.mp3", false},
	}

	for _, f := range files {
		r := q.songInQueue(f.f)
		if r != f.e {
			t.Errorf("songInQueue (%s) was incorrect, got:%v want:%v", f.f, r, f.e)
		}
	}
}

func TestPullNextSong(t *testing.T) {
	q := generateFakeQueue()
	q.currentUser = ""

	tbl := []struct {
		f string
		e bool
	}{
		{"0lol.mp3", true},
		{"2lol.mp3", false},
		{"3foo.mp3", false},
		{"0lol1.mp3", true},
		{"1lol1.mp3", true},
		{"2lol1.mp3", true},
		{"", true},
	}

	for _, v := range tbl {
		next, _ := q.pullNextSong()
		t.Log(q.currentUser)
		if (next.File == v.f) != v.e {
			t.Errorf("pullNextSong was incorrect (%s == %s) got:%v want:%v", next.File, v.f, !v.e, v.e)
		}

	}
}

func TestQueueNextSong(t *testing.T) {
	// q := generateFakeQueue()
	// q.queueNextSong()
}

func TestGetUserQueue(t *testing.T) {
	q := generateFakeQueue()
	e := q.usersQueues[1]
	g := q.getUserQueue(e.user)

	if e != g {
		t.Errorf("GetUserQueue was incorrect got:%s want:%s", g.user, e.user)
	}

	name := "testing"
	g = q.getUserQueue(name)

	if name != g.user {
		t.Errorf("GetUserQueue was incorrect got:%s want:%s", g.user, name)
	}
}
