package mpd

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

type MpdClient struct {
	addr string
	con  *mpd.Client
}

type rankItem struct {
	File string
	Rank int
}

type MpdStatus struct {
	Bitrate        string
	Duration       string
	Elapsed        string
	State          string
	PlaylistLength int
	Song           int
}

func NewMpdClient(addr string) (*MpdClient, error) {
	client, err := mpd.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	c := &MpdClient{
		addr: addr,
		con:  client,
	}

	go func() {
		for range time.Tick(20 * time.Second) {
			client.Ping()
		}
	}()

	return c, nil
}

func (c *MpdClient) searchInLibrary(search string) ([]rankItem, error) {
	songs, err := c.con.ListAllInfo("")
	if err != nil {
		return nil, err
	}

	var items []rankItem

	for _, song := range songs {
		if _, exists := song["file"]; exists {
			rank := fuzzy.RankMatch(
				search,
				song["file"],
			)

			if rank != -1 {
				items = append(items, rankItem{
					File: song["file"],
					Rank: rank,
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
