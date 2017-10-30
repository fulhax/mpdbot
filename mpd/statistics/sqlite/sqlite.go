package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/fulhax/mpdbot/mpd/statistics"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct{ DB *sql.DB }

func New(file string) (statistics.Storage, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	q := `
		CREATE TABLE IF NOT EXISTS songs
		(
			id integer not null primary key,
			title varchar,
			file varchar,
			user varchar,
			createdAt datetime
		)
	`
	_, err = db.Exec(q)
	if err != nil {
		return nil, fmt.Errorf("Error while migrating database: %v", err)
	}

	return &sqlite{db}, nil
}

func (s *sqlite) AddSong(title string, file string, user string) error {
	q := `
		INSERT INTO songs (title, file, user, createdAt)
		VALUES (?, ?, ?, date('NOW'))
	`

	_, err := s.DB.Exec(q, title, file, user)

	return err
}

func (s *sqlite) GetTop(limit int) ([]statistics.SongStats, error) {
	return s.songStatisticsQuery(`
		SELECT title, file, count(file) as timesQueued
		FROM songs
		GROUP BY file
		ORDER BY count(file)
		LIMIT 0,$1
	`, limit)
}

func (s *sqlite) GetUserTop(u string, limit int) ([]statistics.SongStats, error) {
	return s.songStatisticsQuery(`
		SELECT title, file, count(file) as timesQueued
		FROM songs
		WHERE user = ?
		GROUP BY file
		ORDER BY count(file)
		LIMIT 0, ?
	`, u, limit)
}

func (s *sqlite) songStatisticsQuery(q string, args ...interface{}) ([]statistics.SongStats, error) {
	rows, err := s.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []statistics.SongStats

	for rows.Next() {
		var i statistics.SongStats
		rows.Scan(&i.Song.Title, &i.Song.File, &i.TimesQueued)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil

}

func (s *sqlite) Close() error {
	return s.DB.Close()
}
