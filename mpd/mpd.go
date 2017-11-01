package mpd

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

type client interface {
	SearchInLibrary(string) ([]rankItem, error)
	AddSong(string) error
	Play(int) error
	GetStatus() (MpdStatus, error)
	GetState() string
	CurrentSong() (string, error)
	GetRandomSong() (string, error)
	GetAllSongs() ([]mpd.Attrs, error)
	GetPlaylist() ([]mpd.Attrs, error)
	Next() error
	Update() error
	Previous() error
	ClearPlaylist() error
}

type MpdClient struct {
	addr     string
	password string
	con      *mpd.Client
}

type rankItem struct {
	Title string
	File  string
	Rank  int
}

type MpdStatus struct {
	Bitrate        string
	Duration       string
	Elapsed        string
	State          string
	PlaylistLength int
	Song           int
}

func NewMpdClient(addr string, password string) (*MpdClient, error) {
	client, err := mpd.DialAuthenticated("tcp", addr, password)

	if err != nil {
		return nil, err
	}

	c := &MpdClient{
		addr:     addr,
		password: password,
		con:      client,
	}

	go func() {
		for range time.Tick(20 * time.Second) {
			client.Ping()
		}
	}()

	return c, nil
}

func (c *MpdClient) SearchInLibrary(search string) ([]rankItem, error) {
	songs, err := c.GetAllSongs()
	if err != nil {
		return nil, err
	}

	var items []rankItem

	search = strings.ToLower(search)

	for _, song := range songs {
		if _, exists := song["file"]; exists {
			title := ""
			if song["Title"] != "" && song["Artist"] != "" {
				title = fmt.Sprintf("%s - %s", song["Artist"], song["Title"])
			} else {
				title = song["file"]
			}

			rank := fuzzy.RankMatch(
				search,
				strings.ToLower(title),
			)

			if rank != -1 {
				items = append(items, rankItem{
					File:  song["file"],
					Title: title,
					Rank:  rank,
				})
			}
		}
	}

	if len(items) == 0 {
		return nil, nil
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Rank < items[j].Rank
	})

	return items, nil
}

func (c *MpdClient) AddSong(song string) error {
	err := c.con.Add(song)
	if err != nil {
		return err
	}

	return nil
}

func (c *MpdClient) Play(pos int) error {
	err := c.con.Play(pos)
	if err != nil {
		return err
	}

	return nil
}

func (c *MpdClient) GetStatus() (MpdStatus, error) {
	s, err := c.con.Status()
	if err != nil {
		return MpdStatus{}, err
	}

	pLen, _ := strconv.Atoi(s["playlistlength"])
	song, _ := strconv.Atoi(s["song"])

	status := MpdStatus{
		s["bitrate"],
		s["duration"],
		s["elapsed"],
		s["state"],
		pLen,
		song,
	}

	return status, nil
}

func (c *MpdClient) GetState() string {
	return ""
}

func (c *MpdClient) CurrentSong() (string, error) {
	song, err := c.con.CurrentSong()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s - %s", song["Artist"], song["Title"]), nil
}

// GetRandomSong returns random song from mpd library
func (c *MpdClient) GetRandomSong() (string, error) {
	songs, err := c.GetAllSongs()

	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().Unix())
	r := rand.Intn(len(songs) - 1)

	return songs[r]["file"], nil
}

// GetAllSongs returns all songs in MPD library
func (c *MpdClient) GetAllSongs() ([]mpd.Attrs, error) {
	c.con.ListAllInfo("/")
	songs, err := c.con.ListAllInfo("")
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (c *MpdClient) GetPlaylist() ([]mpd.Attrs, error) {
	items, err := c.con.PlaylistInfo(-1, -1)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *MpdClient) Next() error {
	err := c.con.Next()
	if err != nil {
		return err
	}

	return nil
}

func (c *MpdClient) Update() error {
	_, err := c.con.Update("")
	if err != nil {
		return err
	}

	return nil
}

func (c *MpdClient) Previous() error {
	err := c.con.Previous()
	if err != nil {
		return err
	}

	return nil
}

func (c *MpdClient) ClearPlaylist() error {
	status, err := c.GetStatus()
	if err != nil {
		return err
	}

	return c.con.Delete(0, status.PlaylistLength)
}
