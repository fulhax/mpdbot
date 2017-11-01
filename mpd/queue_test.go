package mpd

import (
	"fmt"
	"testing"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/fulhax/mpdbot/mpd/statistics/sqlite"
)

type TestClient struct{}

func (t TestClient) SearchInLibrary(z string) ([]rankItem, error) {
	return []rankItem{}, nil
}
func (t TestClient) AddSong(s string) error            { return nil }
func (t TestClient) Play(z int) error                  { return nil }
func (t TestClient) GetStatus() (MpdStatus, error)     { return MpdStatus{}, nil }
func (t TestClient) GetState() string                  { return "" }
func (t TestClient) CurrentSong() (string, error)      { return "", nil }
func (t TestClient) GetAllSongs() ([]mpd.Attrs, error) { return []mpd.Attrs{}, nil }
func (t TestClient) GetPlaylist() ([]mpd.Attrs, error) { return []mpd.Attrs{}, nil }
func (t TestClient) GetRandomSong() (string, error)    { return "random.mp3", nil }
func (t TestClient) Next() error                       { return nil }
func (t TestClient) Update() error                     { return nil }
func (t TestClient) Previous() error                   { return nil }
func (t TestClient) ClearPlaylist() error              { return nil }

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

	DB, _ := sqlite.New(":memory:")
	q := QueueHandler{
		Client:       &TestClient{},
		currentUser:  "user1",
		usersQueues:  uq,
		StatsStorage: DB,
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
		if (next.File == v.f) != v.e {
			t.Errorf("pullNextSong was incorrect (%s == %s) got:%v want:%v", next.File, v.f, !v.e, v.e)
		}

	}
}

func TestQueueNextSong(t *testing.T) {
	q := generateFakeQueue()
	f := q.queueNextSong()

	e := "0lol1.mp3"
	if f == e {
		t.Errorf("QueueNextSong was incorrect got:%s want:%s", f, e)
	}

	// Set empty queue
	q.usersQueues = []*userQueue{
		{user: "test"},
	}

	f = q.queueNextSong()
	e = "random.mp3"
	if f != e {
		t.Errorf("QueueNextSong was incorrect got:%s want:%s", f, e)
	}

	if "test" == q.currentUser {
		t.Errorf("QueueNextSong was incorrect, expected test got %s", q.currentUser)
	}
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

func TestAddToQueue(t *testing.T) {
	q := generateFakeQueue()

	_, err := q.AddToQueue("user1", "test", "0lol.mp3")

	if nil == err {
		t.Errorf("AddToQueue was incorrect, expected error (already in queue)")
	}

	// Set empty queue
	user := "test"
	file := "test.mp3"
	title := "title"
	q.usersQueues = []*userQueue{
		{user: user},
	}

	i, err := q.AddToQueue(user, title, file)

	if nil != err {
		t.Errorf("AddToQueue was incorrect expected no error: %v", err)
	}

	if i.File != file {
		t.Errorf("AddToQueue was incorrect expected file:%s got %s", file, i.File)
	}

	if i.Title != title {
		t.Errorf("AddToQueue was incorrect expected file:%s got %s", title, i.Title)
	}

	if i.User != user {
		t.Errorf("AddToQueue was incorrect expected file:%s got %s", user, i.User)
	}

	if i != q.usersQueues[0].queue[0] {
		t.Errorf("AddToQueue song not added to user queue")
	}
}
