package statistics

type Storage interface {
	AddSong(string, string, string) error
	GetTop(int) ([]SongStats, error)
	GetUserTop(string, int) ([]SongStats, error)
	Close() error
}

type Song struct {
	Title string
	File  string
}

type SongStats struct {
	Song        Song
	TimesQueued int64
}

type AddedSong struct {
	Song      Song
	User      string
	CreatedAt string
}
