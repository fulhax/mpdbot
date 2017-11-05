package mpd

import (
	"fmt"
	"testing"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/fulhax/mpdbot/mpd/statistics"
	"github.com/fulhax/mpdbot/mpd/statistics/sqlite"
)

type TestClient struct{}

func (t TestClient) SearchInLibrary(z string) ([]rankItem, error) {
	return []rankItem{}, nil
}
func (t TestClient) AddSong(s string) error            { return nil }
func (t TestClient) Play(z int) error                  { return nil }
func (t TestClient) GetStatus() (mpdStatus, error)     { return mpdStatus{}, nil }
func (t TestClient) GetState() string                  { return "" }
func (t TestClient) CurrentSong() (string, error)      { return "", nil }
func (t TestClient) GetAllSongs() ([]mpd.Attrs, error) { return []mpd.Attrs{}, nil }
func (t TestClient) GetPlaylist() ([]mpd.Attrs, error) { return []mpd.Attrs{}, nil }
func (t TestClient) GetRandomSong() (string, error)    { return "random.mp3", nil }
func (t TestClient) Next() error                       { return nil }
func (t TestClient) Update() error                     { return nil }
func (t TestClient) Previous() error                   { return nil }
func (t TestClient) ClearPlaylist() error              { return nil }

func addSongsToStats(u string, c int, stats statistics.Storage) []string {
	var songs []string
	for i := 0; i < c; i++ {
		s := fmt.Sprintf("song%d.mp3", i)
		songs = append(songs, s)
		stats.AddSong(s, s, u)
	}

	return songs
}

func generateFakeQueue() QueueHandler {
	uq := make([]*userQueue, 3)
	for k := range uq {
		u := fmt.Sprintf("user%d", k)
		uq[k] = &userQueue{
			user: u,
			queue: []queueItem{
				queueItem{File: fmt.Sprintf("%dlol.mp3", k), Added: time.Now(), User: u},
				queueItem{File: fmt.Sprintf("%dlol1.mp3", k), Added: time.Now(), User: u},
			},
		}
	}
	s, _ := sqlite.New(":memory:")
	q := QueueHandler{
		Client:       &TestClient{},
		currentUser:  "user1",
		usersQueues:  uq,
		statsStorage: s,
	}

	return q
}

func TestNewQueueHandler(t *testing.T) {
	c := &TestClient{}
	s, _ := sqlite.New(":memory:")
	q := NewQueueHandler(c, s)

	if q.statsStorage != s {
		t.Errorf("stats storage invalid")
	}

	if q.Client != c {
		t.Errorf("wrong client")
	}
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

func TestQueueHandlerPullNextSong(t *testing.T) {
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
			t.Errorf("QueueHandler.pullNextSong was incorrect (%s == %s) got:%v want:%v", next.File, v.f, !v.e, v.e)
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
		t.Errorf("AddToQueue was incorrect expected title:%s got %s", title, i.Title)
	}

	if i.User != user {
		t.Errorf("AddToQueue was incorrect expected user:%s got %s", user, i.User)
	}

	if i != q.usersQueues[0].queue[0] {
		t.Errorf("AddToQueue song not added to user queue")
	}
}

func TestDisableAutoDj(t *testing.T) {
	uq := userQueue{autodj: false}
	uq.DisableAutoDj()
	e := false

	if uq.autodj != e {
		t.Errorf("DisableAutoDj autodj was incorrect expected:%v got:%v", e, uq.autodj)
	}
}

func TestEnableAutoDj(t *testing.T) {
	h := generateFakeQueue()
	uq := userQueue{user: "test", autodj: false, handler: &h}
	err := uq.EnableAutoDj()

	if err == nil {
		t.Errorf("EnableAutoDj expected error")
	}

	addSongsToStats("test", 20, h.statsStorage)

	err = uq.EnableAutoDj()
	if err != nil {
		t.Errorf("EnableAutoDj should not return error: %v", err)
	}

	if uq.autodj != true {
		t.Errorf("EnableAutoDj uq.autodj was incorrect expected: true got: %v", uq.autodj)
	}
}

func TestGetRandomSong(t *testing.T) {
	h := generateFakeQueue()
	uq := userQueue{user: "test", autodj: false, handler: &h}
	err := uq.EnableAutoDj()

	qi, err := uq.getRandomSong()
	if err == nil {
		t.Errorf("getRandomSong should return error if stats is empty")
	}

	songs := addSongsToStats("test", 1, h.statsStorage)

	qi, err = uq.getRandomSong()
	if err != nil {
		t.Errorf("getRandomSong should not return error: %v", err)
	}

	s := songs[0]
	if qi.File != songs[0] {
		t.Errorf("getRandomSong was incorrect expected song:%s got %s", s, qi.File)
	}
}

func TestUserQueuePullNextSongAutoDj(t *testing.T) {
	h := generateFakeQueue()
	uq := userQueue{user: "test", autodj: true, handler: &h}

	songs := addSongsToStats("test", 1, h.statsStorage)

	qi := uq.pullNextSong()
	s := songs[0]
	if qi.File != s {
		t.Errorf("UserQueue.PullNextSongAutoDj was incorrect expected song:%s got %s", s, qi.File)
	}
}

func TestToggleAutodj(t *testing.T) {
	h := generateFakeQueue()
	uq := h.getUserQueue("test")
	addSongsToStats("test", 15, h.statsStorage)

	v, _ := h.ToggleAutodj("test")
	e := true
	if uq.autodj != v {
		t.Errorf("ToggleAutodj was incorrect expected:%v got %v", e, v)
	}

	if v != e {
		t.Errorf("ToggleAutodj was incorrect expected:%v got %v", e, v)
	}

	v, _ = h.ToggleAutodj("test")
	e = false
	if uq.autodj != v {
		t.Errorf("ToggleAutodj was incorrect expected:%v got %v", e, v)
	}

	if v != e {
		t.Errorf("ToggleAutodj was incorrect expected:%v got %v", e, v)
	}
}
